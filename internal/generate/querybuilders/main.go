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
		types.Interface: true,
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

	return typesList, disgord, nil
}

func Exported(name string) bool {
	firstChar := string(name[0])
	return firstChar == strings.ToUpper(firstChar)
}

var disgordTypePrefix = "disgord."

type Entries struct {
	Entries []*TypeWrapper
	Ctx     *Context
}

func (e *Entries) Sort() {
	sort.Slice(e.Entries, func(i, j int) bool {
		name := func(tw *TypeWrapper) string {
			return strings.ToLower(tw.Type.Name.Name)
		}
		return name(e.Entries[i]) < name(e.Entries[j])
	})
}

func (e *Entries) LinkCtx() {
	for i := range e.Entries {
		e.Entries[i].ctx = e.Ctx
	}
}

func (e *Entries) SetPackageName(pkg string) {
	e.Ctx = NewContext(pkg)
	e.LinkCtx()
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

	entries := &Entries{
		Entries: queryBuilders,
	}
	entries.Sort()

	entries.SetPackageName("disgord")
	makeFile(entries, pkg.SourcePath+"/internal/generate/querybuilders/disgord_QueryBuilderNop.gotpl", pkg.SourcePath+"/query_builders_nop_gen.go")

	for _, pkgName := range []string{"disgordutil"} {
		templateFile := fmt.Sprintf("%s/internal/generate/querybuilders/%s_QueryBuilderNop.gotpl", pkg.SourcePath, pkgName)
		destinationFile := fmt.Sprintf("%s/%s/query_builders_nop_gen.go", pkg.SourcePath, pkgName)
		entries.SetPackageName(pkgName)
		makeFile(entries, templateFile, destinationFile)
	}
}

func makeFile(entries *Entries, tplFile, target string) {
	fMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
		"RemovePointer": func(s string) string {
			if s != "" && s[0] == '*' {
				return s[1:]
			}
			return s
		},
	}

	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	if err := tpl.Execute(&b, entries.Entries); err != nil {
		panic(err)
	}

	// Format it according to gofmt standards
	formatted, err := format.Source(b.Bytes())
	if err != nil {
		panic(err)
	}

	//fmt.Println("#####################################")
	//fmt.Println(string(formatted))
	//fmt.Println("#####################################")

	// And write it.
	if err = ioutil.WriteFile(target, formatted, 0644); err != nil {
		panic(err)
	}
}

type TypeWrapper struct {
	hasFlags              bool
	withFlagsReturnType   string
	hasContext            bool
	withContextReturnType string
	*types.Type
	ctx *Context
}

func (t *TypeWrapper) init() {
	for name, m := range t.Methods {
		if name != "WithContext" && name != "WithFlags" {
			continue
		}

		returnType := m.Signature.Results[0].Name.Name
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

func (t *TypeWrapper) DiscordTypePrefix() string {
	if t.ctx.Package == "disgord" {
		return ""
	}
	return "disgord"
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
	return t.ctx.TypeNameWithPackage(t.withContextReturnType, "disgord")
}

func (t *TypeWrapper) WithFlagsReturnType() string {
	return t.ctx.TypeNameWithPackage(t.withContextReturnType, "disgord")
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

		fields = append(fields, &FieldWrapper{&TypeWrapper{Type: m}, name, t.ctx})
	}
	sort.Slice(fields, func(i, j int) bool {
		return strings.ToLower(fields[i].Name) < strings.ToLower(fields[j].Name)
	})
	return fields
}

type FieldWrapper struct {
	Type *TypeWrapper
	Name string
	ctx  *Context
}

func (f *FieldWrapper) TypeName() string {
	return f.Type.TypeName()
}

func (f *FieldWrapper) MethodName() string {
	return f.Name
}

func (f *FieldWrapper) Parameters() string {
	var params []*Type
	for _, p := range f.Type.Signature.Parameters {
		params = append(params, &Type{p, false})
	}
	if len(params) == 0 {
		return ""
	}

	params[len(params)-1].variadic = f.Type.Signature.Variadic
	s := ""
	for _, p := range params {
		if s != "" {
			s += ", "
		}

		name := TypeToLiteral(f.ctx, p)
		s += "_ " + name
	}
	return s
}

func (f *FieldWrapper) ReturnTypes() string {
	var results []*Type
	for _, p := range f.Type.Signature.Results {
		results = append(results, &Type{p, false})
	}
	s := ""
	for _, result := range results {
		if s != "" {
			s += ", "
		}

		s += TypeToLiteral(f.ctx, result)
	}

	if len(f.Type.Signature.Results) > 1 {
		return fmt.Sprintf("(%s)", s)
	}
	return s
}

func (f *FieldWrapper) ReturnValues() string {
	nils := ""
	for _, result := range f.Type.Signature.Results {
		if nils != "" {
			nils += ", "
		}
		nils += ZeroValue(&TypeWrapper{Type: result})
	}
	return nils
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

func ZeroValue(t *TypeWrapper) (v string) {
	if strings.HasSuffix(t.Name.Name, "Snowflake") {
		return "0"
	}

	switch t.Type.Kind {
	case types.Slice:
		v = "nil"
	case types.Pointer, types.Interface:
		v = "nil"
	case types.Struct:
		v = t.TypeName() + "{}"
	case types.Alias:
		v = ZeroValue(&TypeWrapper{Type: t.Type.Underlying})
	case types.Builtin:
		switch t.Type.Name.Name {
		case "bool":
			v = "false"
		case "int", "uint", "int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64":
			v = "0"
		case "float64", "float32":
			v = "0.0"
		case "string":
			v = `""`
		default:
			v = "0 // ++"
		}
	default:
		v = "0  // -"
	}
	return v
}

func CleanName(name string) string {
	replace := func(chars []string) string {
		for _, char := range chars {
			name = strings.Replace(name, char, "", -1)
		}
		return name
	}

	if strings.Contains(name, "func") {
		return name
	}

	return replace([]string{")", "(", "/"})
}

//func MakeTypeNameCompilable(ctx *Context, t *types.Type) string {
//
//	//isDisgordType := false
//	//isDisgordType = isDisgordType || strings.Contains(t.Name.Name, "disgord")
//	//isDisgordType = isDisgordType || strings.Contains(t.Name.Package, "disgord")
//	//isDisgordType = isDisgordType && !strings.Contains(t.Name.Package, "disgord/")
//	//
//	//name := t.Name.Name
//	//name = strings.Replace(name, "github.com/andersfylling/", "", 1)
//	//name = StripTypePrefix(t.Kind, name)
//	//
//	//if t.Kind == types.Slice && strings.Contains(name, "ChannelCreate") {
//	//	fmt.Println(234)
//	//}
//	//
//	//// for reactor handlers we need to identify the func definition
//	//// weird edge case
//	//
//	//if isDisgordType && !strings.Contains(name, disgordTypePrefix) {
//	//	name = disgordTypePrefix + name
//	//}
//	//if strings.Contains(name, "/") {
//	//	joints := strings.Split(name, "/")
//	//	name = joints[len(joints)-1]
//	//} else if !isDisgordType && strings.Contains(t.Name.Package, "/"){
//	//	joints := strings.Split(t.Name.Package, "/")
//	//	name = joints[len(joints)-1] + "." + name
//	//}
//	//
//	//// edge case
//	//if disgordTypePrefix == "" && strings.Contains(name, "disgord.") {
//	//	name = strings.Replace(name, "disgord.", "", 1)
//	//}
//	//
//	//// TODO: improve snowflake..
//	//if strings.HasSuffix(name, "Snowflake") {
//	//	name = disgordTypePrefix + "Snowflake"
//	//}
//	//
//	//name = TypePrefix(t) + name
//	//return CleanName(name)
//
//	return TypeToLiteral(ctx, t)
//}

func TypePrefix(t *types.Type) string {
	switch t.Kind {
	case types.Pointer:
		return "*"
	case types.Slice:
		return "[]"
	case types.Chan:
		return "chan "
	default:
		return ""
	}
}

func StripTypePrefix(kind types.Kind, name string) string {
	switch kind {
	case types.Pointer:
		return strings.Replace(name, "*", "", -1)
	case types.Slice:
		return strings.Replace(name, "[]", "", -1)
	case types.Chan:
		// TODO: support direction
		return strings.Replace(name, "chan", "", -1)
	default:
		return name
	}
}

func IsReactorHandler(t *types.Type) bool {
	if t.Kind == types.Slice || t.Kind == types.Array {
		return IsReactorHandler(t.Elem)
	}

	if t.Kind == types.Func {

	}
	return false
}

func NewContext(Package string) *Context {
	return &Context{
		Imports: map[string]struct{}{},
		Package: Package,
	}
}

type Context struct {
	Imports map[string]struct{}
	Package string
}

func (ctx *Context) TypeName(n types.Name) string {
	segments := strings.Split(n.Package, "/")
	if len(segments) == 0 {
		return n.Name
	}

	packageName := segments[len(segments)-1]
	return ctx.TypeNameWithPackage(n.Name, packageName)
}

func (ctx *Context) TypeNameWithPackage(name string, packageName string) string {
	if packageName == ctx.Package {
		return name
	}
	if packageName == "" {
		return name
	}
	return fmt.Sprintf("%s.%s", packageName, name)
}

type Type struct {
	*types.Type
	variadic bool
}

func (t *Type) Elem() *Type {
	return &Type{t.Type.Elem, false}
}

// TODO: support func as return types (variadic)

func (t *Type) Underlying() *Type {
	return &Type{t.Type.Underlying, false}
}

// TypeToLiteral limited to function declarations. Not entire struct definitions. Simply parameters or return values
func TypeToLiteral(ctx *Context, t *Type) string {
	switch t.Kind {
	case types.Builtin:
		return t.Name.Name
	case types.Alias:
		// edge case
		if t.Name.Name == "Snowflake" {
			return ctx.TypeNameWithPackage(t.Name.Name, "disgord")
		} else {
			return TypeToLiteral(ctx, t.Underlying())
		}
	case types.Struct:
		return ctx.TypeName(t.Name)
	case types.Interface:
		return ctx.TypeName(t.Name)
	case types.Slice:
		var form string
		if t.variadic {
			form = "..."
		} else {
			form = "[]"
		}
		return fmt.Sprintf("%s%s", form, TypeToLiteral(ctx, t.Elem()))
	case types.Pointer:
		return fmt.Sprintf("*%s", TypeToLiteral(ctx, t.Elem()))
	case types.Func:
		return FuncToLiteral(ctx, t.Signature)
	case types.Chan:
		return fmt.Sprintf("chan %s", TypeToLiteral(ctx, t.Elem()))
	}
	panic("type kind is not supported: " + t.Kind)
}

func FuncToLiteral(ctx *Context, t *types.Signature) string {
	var parameters []string
	for _, p := range t.Parameters {
		parameters = append(parameters, TypeToLiteral(ctx, &Type{p, false}))
	}

	var returnVals []string
	for _, r := range t.Results {
		parameters = append(parameters, TypeToLiteral(ctx, &Type{r, false}))
	}

	var returnStmt string
	switch len(returnVals) {
	case 0:
		returnStmt = ""
	case 1:
		returnStmt = returnVals[0]
	default:
		returnStmt = fmt.Sprintf("(%s)", strings.Join(returnVals, ", "))
	}

	var paramStmt string
	paramStmt = strings.Join(parameters, ", ")

	decl := fmt.Sprintf("func(%s) %s", paramStmt, returnStmt)
	return strings.TrimSpace(decl)
}
