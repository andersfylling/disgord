package blankvar

import (
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

	var counter int
	Parse(file, func(a, b *ast.Ident) {
		counter++
	})
	if counter == 0 {
		panic("there should at least be one compile time interface implementation check")
	}
}
