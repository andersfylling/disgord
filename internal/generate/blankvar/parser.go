package blankvar

import (
	"go/ast"
	"go/token"
)

func Parse(file *ast.File, cb func(implementer, contract *ast.Ident)) {
	for _, item := range file.Decls {
		var gdecl *ast.GenDecl
		var ok bool
		if gdecl, ok = item.(*ast.GenDecl); !ok {
			continue
		}

		if gdecl.Tok != token.VAR {
			continue
		}

		specs := item.(*ast.GenDecl).Specs
		for i := range specs {
			vs := specs[i].(*ast.ValueSpec)
			if len(vs.Names) == 0 || vs.Names[0].Name != "_" {
				continue
			}

			var cExpr *ast.CallExpr
			if cExpr, ok = vs.Values[0].(*ast.CallExpr); !ok {
				continue
			}

			var pExpr *ast.ParenExpr
			if pExpr, ok = cExpr.Fun.(*ast.ParenExpr); !ok {
				continue
			}

			var sExpr *ast.StarExpr
			if sExpr, ok = pExpr.X.(*ast.StarExpr); !ok {
				continue
			}

			var id *ast.Ident
			if id, ok = sExpr.X.(*ast.Ident); !ok {
				continue
			}

			var id2 *ast.Ident
			if id2, ok = vs.Type.(*ast.Ident); !ok {
				continue
			}

			cb(id, id2)
		}
	}
}
