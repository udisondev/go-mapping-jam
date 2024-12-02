package mapp

import (
	"fmt"
	"go/ast"
	"strings"
)

type Param struct {
	spec *ast.Field
	imports []Import
}

func (p *Param) Name() string {
	return p.spec.Names[0].Name
}

func (p *Param) Type() (string, string) {
	switch tt := p.spec.Type.(type) {
	case *ast.Ident:
		return "", tt.Name
	case *ast.SelectorExpr:
		return tt.X.(*ast.Ident).Name, tt.Sel.Name

	default:
		panic(fmt.Sprintf("could not extract type from: %T", tt))
	}
}

func (p *Param) Path() string {
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