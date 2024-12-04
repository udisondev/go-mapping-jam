package mapp

import (
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func deepFieldSearch(f Field, fieldFullName string) (Field, bool) {
	if f.FullName() == fieldFullName {
		return f, true
	}

	fields := f.Fields()
	if len(fields) == 0 {
		return Field{}, false
	}

	for _, ff := range fields {
		expF, found := deepFieldSearch(ff, fieldFullName)
		if found {
			return expF, true
		}
	}

	return Field{}, false
}

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