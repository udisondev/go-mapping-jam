package mapp

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Target struct {
	spec *ast.Field
	r Result
}

func (t *Target) Name() string {
	name := t.r.Name()
	if name != "" {
		return name
	}
	return "t"
}

func (t *Target) Fields() []Field {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
		Fset: token.NewFileSet(),
	}
	path := t.r.Path()
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		panic(err)
	}
	pkg := pkgs[0]

	_, name := t.r.Type()
	obj := pkg.Types.Scope().Lookup(name)
	str, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		panic(fmt.Sprintf("'%s.%s' is not a struct!", path, name))
	}

	fields := make([]Field, 0, str.NumFields())
	for i := 0; i < str.NumFields(); i++ {
		fields = append(fields, Field{
			spec:  str.Field(i),
			owner: str,
			path: t.Name(),
		})
	}

	return fields
}