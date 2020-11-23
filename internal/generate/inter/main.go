package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	goparser "go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

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

func RelevantTypes() (typesList []*types.Type, err error) {
	builder := parser.New()
	if err := builder.AddDir(PKGName); err != nil {
		return nil, fmt.Errorf("unable to add disgord package to gengo-parser builder. %w", err)
	}

	universe, err := builder.FindTypes()
	if err != nil {
		return nil, fmt.Errorf("unable to find types for disgord package. %w", err)
	}

	disgord := universe.Package(PKGName)
	for _, typeData := range disgord.Types {
		if accepted, ok := validTypes[typeData.Kind]; !ok || !accepted {
			continue
		}

		// skip unexported types
		// if strings.ToUpper(name[:1]) != name[:1] {
		// 	continue
		// }

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
		if exportedTypes, exportedTypesErr = RelevantTypes(); exportedTypesErr != nil {
			exportedTypesErr = fmt.Errorf("unable to extract exported types: %w", exportedTypesErr)
		}
		wg.Done()
	}()

	// files, getFilesErr := getFiles("/home/anders/dev/disgord/")
	files, getFilesErr := getFiles(".")

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

	makeFile(enforcers["Reseter"], "internal/generate/inter/Reseter.gotpl", "iface_reseter_gen.go")
	makeFile(enforcers["DeepCopier"], "internal/generate/inter/DeepCopier.gotpl", "iface_deepcopier_gen.go")
	makeFile(enforcers["Copier"], "internal/generate/inter/Copier.gotpl", "iface_copier_gen.go")
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
	char := strings.ToLower(t.Name.Name[:1])
	return char
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
	matches := func(kinds ...types.Kind) bool {
		for _, k := range kinds {
			if f.Type.Kind == k {
				return true
			}
		}

		return false
	}

	if matches(types.Slice, types.Array) {
		return false
	}

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

func (f *FieldWrapper) ZeroValue() (v string) {
	switch f.Type.Kind {
	case types.Slice:
		v = "nil"
	case types.Pointer, types.Interface:
		// TODO: check for non-pointers
		v = "nil"
	case types.Struct:
		v = f.TypeName() + "{}"
	case types.Alias:
		v = (&FieldWrapper{Type: &TypeWrapper{f.Type.Underlying, f.Type.typeImplementations}, Name: f.Name}).ZeroValue()
	case types.Builtin:
		switch f.Type.Name.Name {
		case "bool":
			v = "false"
		case "int", "uint", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
			v = "0"
		case "float64", "float32":
			v = "0.0"
		case "string":
			v = ""
		default:
			v = "0 // ++"
		}
	default:
		v = "0  // -"
	}
	return v
}

func (f *FieldWrapper) TypeName() string {
	return f.Type.TypeName()
}

func (f *FieldWrapper) IsSlice() bool {
	return f.Type.Kind == types.Slice
}

func (f *FieldWrapper) IsArray() bool {
	return f.Type.Kind == types.Array
}

func (f *FieldWrapper) MustCopyEach() bool {
	t := f.Type.Elem.Kind
	return (f.IsSlice() || f.IsArray()) && (t == types.Interface || t == types.Pointer)
}

func (f *FieldWrapper) ElemIsPointer() bool {
	t := f.Type.Elem.Kind
	return (f.IsSlice() || f.IsArray()) && t == types.Pointer
}

func (f *FieldWrapper) EventualBuiltin() bool {
	return f.eventual(types.Builtin)
}

func (f *FieldWrapper) EventualInterface() bool {
	return f.eventual(types.Interface)
}

func (f *FieldWrapper) eventual(kind types.Kind) bool {
	var is func(*types.Type) bool
	is = func(t *types.Type) bool {
		if t == nil {
			return false
		} else if t.Kind == kind {
			return true
		}

		return is(t.Elem)
	}
	return is(f.Type.Elem)
}

func (f *FieldWrapper) SliceType() string {
	if f.Type.Type.Kind != types.Slice {
		panic("this is not a slice!")
	}

	// the type definition after the slice prefix: "[]"
	//  "[]" + "*uint64"
	//  "[]" + "uint64"
	//  "[]" + "interface{}"

	var typeData func(*types.Type) string
	typeData = func(t *types.Type) string {
		if t.Kind == types.Slice {
			return "[]" + typeData(t.Elem)
		}

		var v string
		if t.Kind == types.Pointer {
			v += "*"
		}
		name := t.Name.Name
		if strings.Contains(name, "*") {
			name = name[1:]
		}
		if strings.Contains(name, ".") {
			s := strings.Split(name, ".")
			name = s[len(s)-1]
		}

		return v + name
	}

	e := f.Type.Type.Elem
	return typeData(e)
}
