package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	goparser "go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"

	"github.com/andersfylling/disgord/internal/generate/blankvar"
)

const (
	PKGName = "github.com/andersfylling/disgord"
)

var (
	validTypes = map[types.Kind]bool{
		types.Alias:         true,
		types.Array:         true,
		types.Builtin:       true,
		types.Chan:          false,
		types.DeclarationOf: false,
		types.Func:          false,
		types.Interface:     false,
		types.Map:           true,
		types.Pointer:       true,
		types.Protobuf:      false,
		types.Slice:         true,
		types.Struct:        true,
		types.Unknown:       false,
		types.Unsupported:   false,
	}
)

func ExportedTypes() (typesList []*types.Type, err error) {
	builder := parser.New()
	if err := builder.AddDir(PKGName); err != nil {
		return nil, fmt.Errorf("unable to add disgord package to gengo-parser builder. %w", err)
	}

	universe, err := builder.FindTypes()
	if err != nil {
		return nil, fmt.Errorf("unable to find types for disgord package. %w", err)
	}

	disgord := universe.Package(PKGName)
	for name, typeData := range disgord.Types {
		if accepted, ok := validTypes[typeData.Kind]; !ok || !accepted {
			continue
		}

		// skip unexported types
		if strings.ToUpper(name[:1]) != name[:1] {
			continue
		}

		typesList = append(typesList, typeData)
	}

	return typesList, nil
}

func getFiles(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, _ os.FileInfo, err error) error {
		isGoFile := strings.HasSuffix(path, ".go")
		isGenFile := strings.HasSuffix(path, "_gen.go")
		isTestFile := strings.HasSuffix(path, "_test.go")

		var tmpPath string
		if strings.Contains(path, "/disgord") {
			// full path
			tmpPath = path[len(root):]
		} else {
			// relative path
			tmpPath = path
			// TODO: won't work if it's a ../../../ prefix
		}
		isInSubDir := strings.Contains(tmpPath, "/")

		if path == root || !isGoFile || isInSubDir || isGenFile || isTestFile {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func main() {

	var exportedTypes []*types.Type
	var exportedTypesErr error
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if exportedTypes, exportedTypesErr = ExportedTypes(); exportedTypesErr != nil {
			exportedTypesErr = fmt.Errorf("unable to extract exported types: %w", exportedTypesErr)
		}
		wg.Done()
	}()

	files, getFilesErr := getFiles("/home/anders/dev/disgord/")

	typeImplementations := map[string]([]string){}
	for i := range files {
		file, err := goparser.ParseFile(token.NewFileSet(), files[i], nil, 0)
		if err != nil {
			panic(err)
		}

		blankvar.Parse(file, func(implementer, contract *ast.Ident) {
			typeImplementations[implementer.Name] = append(typeImplementations[implementer.Name], contract.Name)
		})
	}

	wg.Wait()
	if getFilesErr != nil {
		panic(fmt.Errorf("unable to fetch files: %w", getFilesErr))
	}
	if exportedTypesErr != nil {
		panic(fmt.Errorf("unable to load all types/contracts. %w", exportedTypesErr))
	}
	if len(typeImplementations) == 0 {
		panic("no enforcers found!")
	}

	enforcers := map[string][]*TypeWrapper{}
	for _, t := range exportedTypes {
		if data, ok := typeImplementations[t.Name.Name]; ok {
			for _, enforcer := range data {
				enforcers[enforcer] = append(enforcers[enforcer], &TypeWrapper{t, typeImplementations})
			}
		}
	}

	makeFile(enforcers["Reseter"], "internal/generate/inter/Reseter.gotpl", "iface2_reseter_gen.go")
}

func makeFile(implementers []*TypeWrapper, tplFile, target string) {
	fMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
	}

	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	if err := tpl.Execute(&b, implementers); err != nil {
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

type TypeWrapper struct {
	*types.Type
	typeImplementations map[string][]string
}

func (t *TypeWrapper) ShortName() string {
	return strings.ToLower(t.Name.Name[:1])
}

func (t *TypeWrapper) TypeName() string {
	if strings.Contains(t.Name.Name, ".") {
		subs := strings.Split(t.Name.Name, ".")
		return subs[len(subs)-1]
	}
	return t.Name.Name
}

func (t *TypeWrapper) Fields() []*FieldWrapper {
	fields := []*FieldWrapper{}
	for _, m := range t.Members {
		fields = append(fields, &FieldWrapper{&TypeWrapper{m.Type, t.typeImplementations}, m.Name})
	}
	return fields
}

type FieldWrapper struct {
	Type *TypeWrapper
	Name string
}

func (f *FieldWrapper) Resetable() bool {
	typeImplementations := f.Type.typeImplementations
	if interfaces, ok := typeImplementations[f.TypeName()]; ok {
		for _, inter := range interfaces {
			if inter == "Reseter" {
				return true
			}
		}
	}

	return false
}

func (f *FieldWrapper) ZeroValue() string {
	return "0"
}

func (f *FieldWrapper) TypeName() string {
	return f.Type.TypeName()
}
