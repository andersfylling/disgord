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
	"sort"
	"strings"
	"text/template"
	"unicode"
)

func main() {
	file, err := parser.ParseFile(token.NewFileSet(), "events.go", nil, 0)
	if err != nil {
		panic(err)
	}

	// Try to find all event structs defined in the file
	var events []*eventName
	var index = map[string]*eventName{}
	for name, item := range file.Scope.Objects {
		// Only continue if we are dealing with a type declaration
		if item.Kind != ast.Typ {
			continue
		}

		// And only if it's a struct definition
		if _, ok := item.Decl.(*ast.TypeSpec).Type.(*ast.StructType); !ok {
			continue
		}

		event := eventName{varName: name}
		events = append(events, &event)
		index[name] = &event
	}

	// Sort them alphabetically instead of the random iteration order from the maps.
	sort.SliceStable(events, func(i, j int) bool {
		return (*events[i]).varName < (*events[j]).varName
	})

	// Next up, read event/events.go to see which events are actual Discord events
	keysFile, err := parser.ParseFile(token.NewFileSet(), "internal/event/events.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	// Read the const key documentation from event/events.go
	for _, item := range keysFile.Decls {
		// Check if this is a GenDecl and if it has at least 1 spec
		genDecl, ok := item.(*ast.GenDecl)
		if !ok || len(genDecl.Specs) == 0 {
			continue
		}
		// Check if it is a ValueSpec and check if it has at least 1 name
		valSpec, ok := genDecl.Specs[0].(*ast.ValueSpec)
		if !ok || len(valSpec.Names) == 0 {
			continue
		}

		name := valSpec.Names[0].Name
		event, ok := index[name]
		if !ok {
			fmt.Fprintf(os.Stderr, "WARNING: event.%s is defined in event/events.go, but we couldn't find the struct!\n", name)
			continue
		}

		doc := genDecl.Doc.Text()
		if doc == "" {
			fmt.Fprintf(os.Stderr, "WARNING: events.%s has no docs! Please write some!\n", name)
		}

		event.Docs = &doc
	}

	for _, event := range events {
		if event.Docs == nil {
			fmt.Fprintf(os.Stderr, "WARNING: %s is defined in events.go, but has no docs in event/events.go!\n", event.varName)
		}
	}

	// And finally pass the event information to different templates to generate some files
	makeFile(events, "generate/events/events.gohtml", "events_gen.go")
	makeFile(events, "generate/events/cache.gohtml", "cache_gen.go")
	makeFile(events, "generate/events/reactor.gotpl", "reactor_gen.go")
}

func makeFile(events []*eventName, tplFile, target string) {
	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	if err := tpl.Execute(&b, events); err != nil {
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

type eventName struct {
	varName string
	Docs    *string
}

func (e eventName) LowerCaseFirst() string {
	return string(unicode.ToLower(rune(e.varName[0]))) + string(e.varName[1:])
}

func (e eventName) String() string {
	return e.varName
}

func (e eventName) IsDiscordEvent() bool {
	return e.Docs != nil
}

func (e eventName) RenderDocs() string {
	if e.Docs == nil {
		return ""
	}

	return "// Evt" + strings.Replace(*e.Docs, "\n", "\n// ", -1)
}
