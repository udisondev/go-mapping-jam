package mapp

import (
	"fmt"
	"go/ast"
	"strings"
)

type Result struct {
	spec *ast.Field
	imports []Import
}


func (p *Result) Name() string {
	if len(p.spec.Names) == 0 {
		return ""
	}

	return p.spec.Names[0].Name
}

func (p *Result) Type() (string, string) {
	switch tt := p.spec.Type.(type) {
	case *ast.Ident:
		return "", tt.Name
	case *ast.SelectorExpr:
		return tt.X.(*ast.Ident).Name, tt.Sel.Name

	default:
		panic(fmt.Sprintf("could not extract type from: %T", tt))
	}
}

func (p *Result) Path() string {
	alias, _ := p.Type()
	if alias == "" {
		return ""
	}

	for _, i := range p.imports {
		if i.Alias() == alias {
			return strings.ReplaceAll(i.Path(), "\"", "")
		}
	}

	return ""
}