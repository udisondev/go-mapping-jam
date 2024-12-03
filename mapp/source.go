package mapp

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Source struct {
	spec *ast.Field
	Param
}

func (s Source) Fields() []Field {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
		Fset: token.NewFileSet(),
	}
	path := s.Path()
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		panic(err)
	}
	pkg := pkgs[0]

	_, name := s.Type()
	obj := pkg.Types.Scope().Lookup(name)
	str, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		panic(fmt.Sprintf("'%s.%s' is not a struct!", path, name))
	}

	fields := make([]Field, 0, str.NumFields())
	for i := 0; i < str.NumFields(); i++ {
		fields = append(fields, Field{
			spec:      str.Field(i),
			owner:     str,
			fieldPath: s.Name(),
		})
	}

	return fields
}

func (s Source) FieldByFullName(fullName string) (Field, bool) {
	fields := s.Fields()
	for _, f := range fields {	
		expF, found := deepFieldSearch(f, fullName)
		if found {
			return expF, found
		}
	}
	return Field{}, false
}
