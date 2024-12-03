package mapp

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func extractFieldsFromStruct(filedPath, typePath, typeName string) []Field {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
		Fset: token.NewFileSet(),
	}
	pkgs, err := packages.Load(cfg, typePath)
	if err != nil {
		panic(err)
	}
	pkg := pkgs[0]

	obj := pkg.Types.Scope().Lookup(typeName)
	str, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	fields := make([]Field, 0, str.NumFields())
	for i := 0; i < str.NumFields(); i++ {
		fields = append(fields, Field{
			spec:      str.Field(i),
			owner:     str,
			fieldPath: filedPath,
		})
	}

	return fields
}
