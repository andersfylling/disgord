package querybuilderparser

import (
	"go/ast"
	"go/token"
	"strings"
)

func Exported(name string) bool {
	firstChar := string(name[0])
	return firstChar == strings.ToUpper(firstChar)
}

func Parse(file *ast.File, cb func(i *ast.TypeSpec)) {
	for _, item := range file.Decls {
		var gdecl *ast.GenDecl
		var ok bool
		if gdecl, ok = item.(*ast.GenDecl); !ok {
			continue
		}

		if gdecl.Tok != token.TYPE {
			continue
		}

		specs := item.(*ast.GenDecl).Specs
		for i := range specs {
			ts := specs[i].(*ast.TypeSpec)
			name := ts.Name.String()
			if len(name) == 0 {
				continue
			}
			if !strings.HasSuffix(name, "QueryBuilder") {
				continue
			}
			if !Exported(name) {
				continue
			}

			cb(ts)
		}
	}
}
