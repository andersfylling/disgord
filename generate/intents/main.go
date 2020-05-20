package main

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"sort"
	"text/template"
)

func main() {
	file, err := parser.ParseFile(token.NewFileSet(), "internal/gateway/intents.go", nil, 0)
	if err != nil {
		panic(err)
	}

	// Find all intent values
	var intents []string
	for name, item := range file.Scope.Objects {
		if item.Kind != ast.Con {
			continue
		}
		intents = append(intents, name)
	}

	// Sort them alphabetically instead of the random iteration order from the maps.
	sort.SliceStable(intents, func(i, j int) bool {
		return intents[i] < intents[j]
	})

	// And finally pass the event information to different templates to generate some files
	makeFile(intents, "generate/intents/intents.gohtml", "intents_gen.go")
}

func makeFile(intents []string, tplFile, target string) {
	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	if err := tpl.Execute(&b, intents); err != nil {
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
