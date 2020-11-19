package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func main() {
	files, err := getFiles(".")
	if err != nil {
		panic(err)
	}

	structs := make([]*Struct, 0, 150)
	for i := range files {
		structs = append(structs, getJSONStructs(files[i])...)
	}

	// TODO: now that structs holds all json structs in Disgord
	// it's time to use ffjson or another tool to generate the marshal/unmarshal methods

	for _ = range structs {
		//fmt.Println(structs[i].Name)
	}
}

type Struct struct {
	Name string
	Obj  *ast.StructType
}

func getJSONStructs(filename string) (structs []*Struct) {
	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		panic(err)
	}
	for _, item := range file.Decls {
		var gdecl *ast.GenDecl
		var ok bool
		if gdecl, ok = item.(*ast.GenDecl); !ok {
			continue
		}

		if gdecl.Tok != token.TYPE {
			continue
		}

		specs := item.(*ast.GenDecl).Specs
		for i := range specs {
			ts := specs[i].(*ast.TypeSpec)
			if ts.Name == nil || !ts.Name.IsExported() {
				continue
			}

			var st *ast.StructType
			if st, ok = ts.Type.(*ast.StructType); !ok || st.Fields == nil {
				continue
			}

			var hasJSONTags bool
			for j := range st.Fields.List {
				if len(st.Fields.List[j].Names) == 0 {
					continue
				}

				field := st.Fields.List[j]
				if field.Tag != nil && len(field.Tag.Value) > 2 {
					tagStruct := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]) // rm ` wraps
					hasJSONTags = tagStruct.Get("json") != ""
					if hasJSONTags {
						break
					}
				}
			}

			if !hasJSONTags {
				continue
			}

			structs = append(structs, &Struct{
				Name: ts.Name.Name,
				Obj:  st,
			})
		}
	}

	return structs
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
		isTestFile := strings.HasSuffix(results[i], "_test.go")
		isGenFile := strings.HasSuffix(results[i], "_gen.go")
		if results[i] == path || !isGoFile || isInSubDir || isTestFile || isGenFile {
			continue
		}

		files = append(files, results[i])
	}

	return files, nil
}

//
//func makeFile(builders []*builder, tplFile, target string) {
//	fMap := template.FuncMap{
//		"ToUpper":      strings.ToUpper,
//		"ToLower":      strings.ToLower,
//		"Capitalize":   Capitalize,
//		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
//	}
//
//	// Open & parse our template
//	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))
//
//	// Execute the template, inserting all the event information
//	var b bytes.Buffer
//	if err := tpl.Execute(&b, builders); err != nil {
//		panic(err)
//	}
//
//	// Format it according to gofmt standards
//	formatted, err := format.Source(b.Bytes())
//	if err != nil {
//		panic(err)
//	}
//
//	// And write it.
//	if err = ioutil.WriteFile(target, formatted, 0644); err != nil {
//		panic(err)
//	}
//}
