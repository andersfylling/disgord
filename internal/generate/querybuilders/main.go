package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"path"
	"sort"
	"strings"
	"text/template"

	"k8s.io/gengo/parser"
	"k8s.io/gengo/types"
)

const (
	PKGName = "github.com/andersfylling/disgord"
)

var (
	validTypes = map[types.Kind]bool{
		types.Interface:     true,
	}
)

func DisgordTypes() (typesList []*types.Type, p *types.Package, err error) {
	builder := parser.New()
	if err := builder.AddDir(PKGName); err != nil {
		return nil, nil, fmt.Errorf("unable to add disgord package to gengo-parser builder. %w", err)
	}

	universe, err := builder.FindTypes()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to find types for disgord package. %w", err)
	}

	disgord := universe.Package(PKGName)
	for _, typeData := range disgord.Types {
		if accepted, ok := validTypes[typeData.Kind]; !ok || !accepted {
			continue
		}

		typesList = append(typesList, typeData)
	}

	return typesList, disgord,nil
}

func Exported(name string) bool {
	firstChar := string(name[0])
	return firstChar == strings.ToUpper(firstChar)
}

func main() {
	disgordTypes, pkg, err := DisgordTypes()
	if err != nil {
		panic(err)
	}

	var queryBuilders []*TypeWrapper
	for _, t := range disgordTypes {
		name := t.Name.Name
		if !strings.HasSuffix(name, "QueryBuilder") {
			continue
		}
		if !Exported(name) {
			continue
		}

		wrap := &TypeWrapper{Type: t}
		wrap.init()
		queryBuilders = append(queryBuilders, wrap)
	}

	makeFile(queryBuilders, pkg.SourcePath+"/internal/generate/querybuilders/disgordutil_QueryBuilderNop.gotpl", pkg.SourcePath+"/disgordutil/query_builders_nop_gen.go")
}

func makeFile(implementers []*TypeWrapper, tplFile, target string) {
	fMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
	}

	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))

	// sort the enforcers so that the generated output stays the same every time
	sort.Slice(implementers, func(i, j int) bool {
		name := func(tw *TypeWrapper) string {
			return strings.ToLower(tw.Type.Name.Name)
		}
		return name(implementers[i]) < name(implementers[j])
	})

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
	hasFlags bool
	withFlagsReturnType string
	hasContext bool
	withContextReturnType string
	*types.Type
}

func (t *TypeWrapper) init() {
	for name, m := range t.Methods {
		if name != "WithContext" && name != "WithFlags" {
			continue
		}

		returnType := "disgord." + m.Signature.Results[0].Name.Name
		if name == "WithContext" {
			t.hasContext = true
			t.withContextReturnType = returnType
			continue
		}
		if name == "WithFlags" {
			t.hasFlags = true
			t.withFlagsReturnType = returnType
		}
	}
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

func (t *TypeWrapper) WithContextReturnType() string {
	return t.withContextReturnType
}

func (t *TypeWrapper) WithFlagsReturnType() string {
	return t.withFlagsReturnType
}

func (t *TypeWrapper) HasWithContext() bool {
	return t.hasContext
}

func (t *TypeWrapper) HasWithFlags() bool {
	return t.hasFlags
}

func (t *TypeWrapper) Fields() []*FieldWrapper {
	var fields []*FieldWrapper
	for name, m := range t.Methods {
		if name == "WithContext" || name == "WithFlags" {
			continue
		}

		if name == "CreateInvite" {
			fmt.Sprintln("asdasd")
		}
		fields = append(fields, &FieldWrapper{&TypeWrapper{Type: m}, name})
	}
	sort.Slice(fields, func(i, j int) bool {
		return strings.ToLower(fields[i].Name) < strings.ToLower(fields[j].Name)
	})
	return fields
}

type FieldWrapper struct {
	Type *TypeWrapper
	Name string
}

func (f *FieldWrapper) TypeName() string {
	return f.Type.TypeName()
}

func (f *FieldWrapper) MethodName() string {
	return f.Name
}

func (f *FieldWrapper) Parameters() string {
	return ""
}

func (f *FieldWrapper) ReturnTypes() string {
	s := ""
	for _, result := range f.Type.Signature.Results {
		if s != "" {
			s += ", "
		}
		//if result.Kind == types.Pointer {
		//	s += "*"
		//}
		//if result.Kind != types.Builtin {
		//	s += "disgord."
		//}

		name := result.Name.Name
		name = strings.Replace(name, "github.com/andersfylling/", "", 1)
		s += name
	}

	return s
}

func (f *FieldWrapper) ReturnValues() string {
	return "nil"
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
