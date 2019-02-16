package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
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

type builder struct {
	name      string
	fieldName string
}

func (e builder) String() string {
	return e.name
}

func (e builder) FieldName() string {
	return e.fieldName
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

		// must contain Builder in it's name / suffix
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
					panic("" + filename + "#" + name + " must specify a field name for the embedded struct " + RESTBuilder)
				}
				fieldName = field.Names[0].Name
				isRESTBuilder = true
				break
			}

		}
		if !isRESTBuilder {
			continue
		}

		builders = append(builders, &builder{
			name:      name,
			fieldName: fieldName,
		})
	}

	sort.Slice(builders, func(i, j int) bool {
		return builders[i].name < builders[j].name
	})

	return builders
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
	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).ParseFiles(tplFile))

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
