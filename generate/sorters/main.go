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

	e := &env{}

	for i := range files {
		e.Types = append(e.Types, getTypes(files[i])...)
	}

	e.Sorters = []Sorter{
		{
			Field:      "ID", // TODO: check if the field type is correct. regression.
			Ascending:  func(s string) string { return "<" + s },
			Descending: func(s string) string { return ">" + s },
		},
		{
			Field:      "Name",
			Ascending:  func(s string) string { return "<" + s },
			Descending: func(s string) string { return ">" + s },
		},
	}

	for i := range e.Sorters {
		for j := range e.Types {
			if e.Types[j].HasFieldName(e.Sorters[i].Field) {
				e.Sorters[i].Types = append(e.Sorters[i].Types, e.Types[j])
			}
		}
	}

	makeFile(e, "generate/sorters/sorters.gotpl", "sort_gen.go")
}

type Type struct {
	Name   string
	Fields []string
}

func (t *Type) HasFieldName(name string) bool {
	for i := range t.Fields {
		if t.Fields[i] == name {
			return true
		}
	}

	return false
}

type Sorter struct {
	Type       string
	Field      string
	Ascending  func(s string) string // s is whatever to compare against: a[i].ID < b[j].ID => b[j].ID is s
	Descending func(s string) string // s is whatever to compare against: a[i].ID > b[j].ID => b[j].ID is s
	Types      []Type
}

type env struct {
	Sorters []Sorter
	Types   []Type
}

func getTypes(filename string) (types []Type) {
	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		panic(err)
	}

	for name, item := range file.Scope.Objects {
		// Only continue if we are dealing with a type declaration
		if item.Kind != ast.Typ {
			continue
		}

		// And only if it's a struct definition
		typeDecl := item.Decl.(*ast.TypeSpec)
		var structDecl *ast.StructType
		var ok bool
		if structDecl, ok = typeDecl.Type.(*ast.StructType); !ok {
			continue
		}

		t := Type{Name: name}
		for _, field := range structDecl.Fields.List {
			if len(field.Names) == 0 {
				continue
			}
			t.Fields = append(t.Fields, field.Names[0].Name)
		}

		types = append(types, t)
	}

	sort.Slice(types, func(i, j int) bool {
		return types[i].Name < types[j].Name
	})

	return types
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

func makeFile(e *env, tplFile, target string) {
	fMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"Capitalize":   Capitalize,
		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
		"Ascending": func(field, name string) string {
			for i := range e.Sorters {
				if e.Sorters[i].Field == field {
					return e.Sorters[i].Ascending(name + "[j]." + field)
				}
			}
			return ""
		},
		"Descending": func(field, name string) string {
			for i := range e.Sorters {
				if e.Sorters[i].Field == field {
					return e.Sorters[i].Descending(name + "[j]." + field)
				}
			}
			return ""
		},
	}

	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	if err := tpl.Execute(&b, e); err != nil {
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
