package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

func main() {
	files, err := getFiles(".")
	if err != nil {
		panic(err)
	}

	var builders []*builder
	for i := range files {
		builders = append(builders, getAllRESTBuilders(files[i])...)
	}

	makeFile(builders, "generate/restbuilders/methods.gotpl", "restbuilders_gen.go")
}

type condition struct {
	min int
	max int
}

var _ fmt.Stringer = (*condition)(nil)

func (c *condition) String() string {
	return fmt.Sprintf("%d<N<%d", c.min, c.max)
}

func (c *condition) Render(name string) string {
	min := strconv.Itoa(c.min)
	max := strconv.Itoa(c.max)
	// min < b && b < max
	return "(" + min + "<" + name + " && " + name + "<" + max + `, "` + name + ` must be in the range of (` + min + ", " + max + `)")`
}

func GetCondition(param string) *condition {
	if param == "" {
		return nil
	}
	if strings.HasPrefix(param, "(") {
		param = param[1:]
	}
	if strings.HasSuffix(param, ")") {
		param = param[:len(param)-1]
	}

	vals := strings.Split(param, "<")
	if len(vals) != 3 {
		panic(fmt.Sprintf("expects to have 3 values, got: %+v", vals))
	}
	if vals[1] != "N" && vals[1] != "n" {
		panic(fmt.Sprintf("expected middle value to be N, got %s", vals[1]))
	}

	min, err := strconv.Atoi(vals[0])
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(vals[2])
	if err != nil {
		panic(err)
	}

	return &condition{
		min: min,
		max: max,
	}
}

func ProcessValueParam(param string) (name string, cond *condition) {
	if param == "" {
		return "", nil
	}

	if strings.Contains(param, "(") {
		condStr := param[strings.Index(param, "("):]
		cond = GetCondition(condStr)
		name = param[:len(param)-len(condStr)]
	} else {
		name = param
	}

	return name, cond
}

type methodInfo struct {
	Name       string
	MethodName string
	Type       string
	IsSlice    bool
	Cond       *condition
}

func (m *methodInfo) HasCond() bool {
	return m.Cond != nil
}

type builder struct {
	Name      string
	FieldName string
	pos       token.Pos
	Params    []methodInfo
	BasicExec *methodInfo
}

func (b *builder) BasicExecMethod() bool {
	return b.BasicExec != nil
}

func getAllRESTBuilders(filename string) (builders []*builder) {
	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		panic(err)
	}

	// Try to find all event structs defined in the file
	const RESTBuilder = "RESTBuilder"
	for name, item := range file.Scope.Objects {
		// Only continue if we are dealing with a type declaration
		if item.Kind != ast.Typ {
			continue
		}

		// must contain Builder in it's Name / suffix
		if !strings.HasSuffix(name, "Builder") || name == RESTBuilder {
			continue
		}

		// And only if it's a struct definition
		typeDecl := item.Decl.(*ast.TypeSpec)
		var structDecl *ast.StructType
		var ok bool
		if structDecl, ok = typeDecl.Type.(*ast.StructType); !ok {
			continue
		}

		// and if the struct embeds restBuilder
		fields := structDecl.Fields.List
		var isRESTBuilder bool
		var fieldName string
		for _, field := range fields {
			if id, ok := field.Type.(*ast.Ident); ok && id.Name == RESTBuilder {
				if len(field.Names) == 0 {
					panic("" + filename + "#" + name + " must specify a field Name for the embedded struct " + RESTBuilder)
				}
				fieldName = field.Names[0].Name
				isRESTBuilder = true
				break
			}

		}
		if !isRESTBuilder {
			continue
		}

		pos := item.Pos()
		if pos > 300 {
			pos -= 300 // magic
		} else {
			pos = 0
		}
		builders = append(builders, &builder{
			Name:      name,
			FieldName: fieldName,
			pos:       pos,
		})
	}

	sort.Slice(builders, func(i, j int) bool {
		return builders[i].Name < builders[j].Name
	})

	// Read the comment to check for generate instructions
	fileComments, err := parser.ParseFile(token.NewFileSet(), filename, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	const genPrefix = "//generate-rest-params: "
	const genPrefix2 = "//generate-rest-basic-execute: "
	var ok bool
	var genDecl *ast.GenDecl
	for _, item := range fileComments.Decls {
		if len(builders) == 0 {
			break
		}

		// magic
		//if item.Pos() < builders[0].pos {
		//	continue
		//}

		genDecl, ok = item.(*ast.GenDecl)
		if !ok || len(genDecl.Specs) == 0 || genDecl.Doc == nil || len(genDecl.Doc.List) == 0 {
			continue
		}

		spec, ok := genDecl.Specs[0].(*ast.TypeSpec)
		if !ok || spec.Name.IsExported() {
			continue
		}

		// find slice index
		structName := spec.Name.Name
		index := -1
		for i := range builders {
			if builders[i].Name == structName {
				index = i
				break
			}
		}
		if index == -1 {
			continue
		}

		for i := range genDecl.Doc.List {
			comment := genDecl.Doc.List[i].Text
			if !strings.HasPrefix(comment, genPrefix) && !strings.HasPrefix(comment, genPrefix2) {
				continue
			}

			if strings.HasPrefix(comment, genPrefix) {
				var start = len(genPrefix)
				var end int
				if strings.HasSuffix(comment, ",") {
					end = len(comment) - 1
				} else {
					end = len(comment)
				}
				paramsStr := comment[start:end]
				params := strings.Split(paramsStr, ", ")

				tuples := make([]methodInfo, 0, len(params))
				for j := range params {
					param := strings.Split(params[j], ":")
					name := param[0]
					typ, cond := ProcessValueParam(param[1])
					methodName := jsonNameToMethodName(name)
					tuples = append(tuples, methodInfo{
						MethodName: methodName,
						Name:       name,
						Type:       typ,
						Cond:       cond,
					})
				}
				builders[index].Params = tuples
			} else if strings.HasPrefix(comment, genPrefix2) {
				var start = len(genPrefix2)
				var end int
				if strings.HasSuffix(comment, ",") {
					end = len(comment) - 1
				} else {
					end = len(comment)
				}
				paramsStr := comment[start:end]
				params := strings.Split(paramsStr, ", ")
				param := strings.Split(params[0], ":")
				isSlice := strings.HasPrefix(param[1], "[]")
				// let panic if missing params
				builders[index].BasicExec = &methodInfo{
					Name:    param[0],
					Type:    param[1],
					IsSlice: isSlice,
				}
			}

		}
	}

	return builders
}

func jsonNameToMethodName(name string) string {
	words := strings.Split(name, "_")
	var methodName string
	for i := range words {
		if words[i] == "id" {
			methodName += "ID"
		} else {
			methodName += Capitalize(words[i])
		}
	}
	return methodName
}

func Capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
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

func makeFile(builders []*builder, tplFile, target string) {
	fMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"Capitalize":   Capitalize,
		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
	}

	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	if err := tpl.Execute(&b, builders); err != nil {
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
