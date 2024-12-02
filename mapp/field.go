//go:generate go-enum
package mapp

import (
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// FieldType ENUM(
// Basic,
// Struct,
// Slice,
// PointerToBasic,
// PointerToStruct,
// PointerToSlice,
// )
type FieldType uint8

type Field struct {
	spec  *types.Var
	owner *types.Struct
	path  string
}

func (f *Field) Name() string {
	return f.spec.Origin().Name()
}

func (f *Field) FullName() string {
	return f.path + "." + f.Name()
}

func (f *Field) Fields() []Field {
	switch ft := f.spec.Type().(type) {
	case *types.Basic:
		return nil
	case *types.Named:
		_, isStruct := ft.Underlying().(*types.Struct)
		if !isStruct {
			return nil
		}

		path := ft.Obj().Pkg().Path()
		splitedType := strings.Split(f.spec.Origin().Type().String(), ".")
		name := splitedType[len(splitedType)-1]
		return f.returnStructFields(path, name)
	case *types.Pointer:
		switch ft.Underlying().(type) {
		case *types.Struct:
			f.returnStructFields(f.spec.Pkg().Path(), f.spec.Type().String())
		}
	}

	return nil
}

func (f *Field) returnStructFields(path, typeName string) []Field {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
		Fset: token.NewFileSet(),
	}
	pkgs, err := packages.Load(cfg, path)
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
			spec:  str.Field(i),
			owner: str,
			path:  f.FullName(),
		})
	}

	return fields
}
