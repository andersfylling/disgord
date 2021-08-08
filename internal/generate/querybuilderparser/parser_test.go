package querybuilderparser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestParse(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "../../../channel.go", nil, 0)
	if err != nil {
		panic(err)
	}

	var interfaces []string
	var counter int
	Parse(file, func(i *ast.TypeSpec) {
		counter++
		interfaces = append(interfaces, i.Name.String())
	})
	if counter == 0 {
		panic("there should at least be one compile time interface implementation check")
	}

	fmt.Println(interfaces)
}
