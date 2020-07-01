package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
)

type tagInfo struct {
	Name         string
	Omitempty    bool
	ZeroValCheck string // " == 0", ".IsZero()", etc.
}

type fieldInfo struct {
	Name             string
	ZeroVal          string
	typ              string
	Tag              *tagInfo
	resetableStructs *[]structInfo
}

func (f *fieldInfo) HasTag() bool {
	return f.Tag != nil
}

func (f *fieldInfo) Resetable() bool {
	structs := *f.resetableStructs
	name := f.typ
	if strings.Contains(name, " ") {
		// pointer. eg. &{405 PartialEmoji}
		s := strings.Split(name, " ")
		name = s[1][:len(s[1])-1] // remove } suffix
	}
	for i := range structs {
		if name == structs[i].Name {
			return true
		}
	}

	return false
}

type structInfo struct {
	Name      string
	ShortName string
	Fields    []fieldInfo
}

type Enforcer struct {
	Name    string
	Structs []structInfo
}

func typeOfInterest(enfs []Enforcer, name string) bool {
	for i := range enfs {
		for j := range enfs[i].Structs {
			if enfs[i].Structs[j].Name == name {
				return true
			}
		}
	}

	return false
}

func removeStruct(enfs []Enforcer, name string) {
	for i := range enfs {
		for j := range enfs[i].Structs {
			if enfs[i].Structs[j].Name == name {
				enfs[i].Structs = append(enfs[i].Structs[:j], enfs[i].Structs[j+1:]...)
				break
			}
		}
	}
}

func main() {
	files, err := getFiles(".")
	if err != nil {
		panic(err)
	}

	enforcers := []Enforcer{
		{Name: "Reseter"},
		{Name: "URLQueryStringer"},

		{Name: "internalUpdater"},
		{Name: "internalClientUpdater"},
	}
	for i := range files {
		file, err := parser.ParseFile(token.NewFileSet(), files[i], nil, 0)
		if err != nil {
			panic(err)
		}

		addEnforcers(enforcers, file)
	}
	for i := range files {
		file, err := parser.ParseFile(token.NewFileSet(), files[i], nil, 0)
		if err != nil {
			panic(err)
		}

		addStructs(enforcers, file)
	}

	makeFile(enforcers, "generate/interfaces/Reseter.gotpl", "iface_reseter_gen.go")
	makeFile(enforcers, "generate/interfaces/URLQueryStringer.gotpl", "iface_urlquerystringer_gen.go")
	makeFile(enforcers, "generate/interfaces/internalUpdaters.gotpl", "iface_internalupdaters_gen.go")
}

func addStructs(enforcers []Enforcer, file *ast.File) {
	for _, item := range file.Decls {
		var gdecl *ast.GenDecl
		var ok bool
		if gdecl, ok = item.(*ast.GenDecl); !ok {
			continue
		}

		if gdecl.Tok != token.TYPE {
			continue
		}

		var resetables *[]structInfo
		for i := range enforcers {
			if enforcers[i].Name == "Reseter" {
				resetables = &enforcers[i].Structs
			}
		}

		specs := item.(*ast.GenDecl).Specs
		for i := range specs {
			ts := specs[i].(*ast.TypeSpec)
			if ts.Name == nil || !typeOfInterest(enforcers, ts.Name.Name) {
				continue
			}

			if st, ok := ts.Type.(*ast.StructType); ok && st.Fields != nil {
				for j := range st.Fields.List {
					if len(st.Fields.List[j].Names) == 0 {
						continue
					}
					field := st.Fields.List[j]
					name := field.Names[0].Name
					typ := fmt.Sprint(field.Type)

					var tag *tagInfo
					if field.Tag != nil && len(field.Tag.Value) > 2 {
						tagStruct := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]) // rm ` wraps
						detailsStr := tagStruct.Get("urlparam")
						details := strings.Split(detailsStr, ",")
						if len(details) > 0 {
							tag = &tagInfo{
								Name:      details[0],
								Omitempty: strSliceContains(details, "omitempty"),
							}

							if zVal, ok := getZeroVal(typ); ok && tag.Omitempty {
								tag.ZeroValCheck = " == " + zVal
							} else {
								tag.ZeroValCheck = " OMFG"
							}
						}

						if tag.Name == "-" {
							continue
						}
					}

					var zeroInit string
					var ok bool
					if zeroInit, ok = getZeroVal(typ); !ok {
						zeroInit = getZeroInit(typ)
					}
					// fmt.Println(name, " = ", typ, " => ", zeroVal)

					for a := range enforcers {
						for b := range enforcers[a].Structs {
							if enforcers[a].Structs[b].Name == ts.Name.Name {
								info := fieldInfo{Name: name, ZeroVal: zeroInit, Tag: tag, typ: typ, resetableStructs: resetables}
								enforcers[a].Structs[b].Fields = append(enforcers[a].Structs[b].Fields, info)
								break
							}
						}
					}
				}
			} else {
				removeStruct(enforcers, ts.Name.Name)
			}

		}
	}
}

func addEnforcers(enforcers []Enforcer, file *ast.File) {
	for _, item := range file.Decls {
		var gdecl *ast.GenDecl
		var ok bool
		if gdecl, ok = item.(*ast.GenDecl); !ok {
			continue
		}

		if gdecl.Tok != token.VAR {
			continue
		}

		specs := item.(*ast.GenDecl).Specs
		for i := range specs {
			vs := specs[i].(*ast.ValueSpec)
			if len(vs.Names) == 0 || vs.Names[0].Name != "_" {
				continue
			}

			var cExpr *ast.CallExpr
			if cExpr, ok = vs.Values[0].(*ast.CallExpr); !ok {
				continue
			}

			var pExpr *ast.ParenExpr
			if pExpr, ok = cExpr.Fun.(*ast.ParenExpr); !ok {
				continue
			}

			var sExpr *ast.StarExpr
			if sExpr, ok = pExpr.X.(*ast.StarExpr); !ok {
				continue
			}

			var id *ast.Ident
			if id, ok = sExpr.X.(*ast.Ident); !ok {
				continue
			}

			var id2 *ast.Ident
			if id2, ok = vs.Type.(*ast.Ident); !ok {
				continue
			}

			for j := range enforcers {
				if enforcers[j].Name == id2.Name {
					s := structInfo{
						Name:      id.Name,
						ShortName: strings.ToLower(id.Name[:1]),
					}
					enforcers[j].Structs = append(enforcers[j].Structs, s)
					break
				}
			}
		}
	}
}

func strSliceContains(s []string, needle string) bool {
	for i := range s {
		if s[i] == needle {
			return true
		}
	}

	return false
}

func getZeroVal(s string) (result string, success bool) {
	switch s {
	case "int", "int8", "int16", "int32", "int64":
		result = "0"
	case "uint", "uint8", "uint16", "uint32", "uint64":
		result = "0"
	case "float32", "float64":
		result = "0"
	case "Snowflake", "snowflake.ID", "snowflake.Snowflake", "depalias.Snowflake", "MessageType", "MessageFlag":
		result = "0"
	case "string":
		result = ""
		success = true
	case "bool":
		result = "false"
	case "nil":
		result = s
		// TODO: find out what the original data type is
	case "VerificationLvl", "DefaultMessageNotificationLvl", "ExplicitContentFilterLvl", "MFALvl", "Discriminator", "PremiumType", "PermissionBit", "activityFlag", "acitivityType":
		result = "0"
	}

	if !success && result != "" {
		success = true
	}

	return result, success
}

func getZeroInit(s string) string {
	switch s {
	case "time.Time", "&{time Time}":
		return "time.Unix(0, 0)"
	case "Timestamp":
		return "Timestamp(time.Unix(0, 0))"
	case "UserFlag":
		return "0"
	}

	if strings.HasPrefix(s, "&{") {
		return "nil"
	}

	return s + "{}"
}

func getFiles(path string) (files []string, err error) {
	var results []string
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		results = append(results, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	for i := range results {
		isGoFile := strings.HasSuffix(results[i], ".go")
		isInSubDir := strings.Contains(results[i], "/")
		isGenFile := strings.HasSuffix(results[i], "_gen.go")
		if results[i] == path || !isGoFile || isInSubDir || isGenFile {
			continue
		}

		files = append(files, results[i])
	}

	return files, nil
}

func makeFile(enforcers []Enforcer, tplFile, target string) {
	fMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
	}

	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	if err := tpl.Execute(&b, enforcers); err != nil {
		panic(err)
	}

	// Format it according to gofmt standards
	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}

	// And write it.
	if err = ioutil.WriteFile(target, formatted, 0644); err != nil {
		panic(err)
	}
}
