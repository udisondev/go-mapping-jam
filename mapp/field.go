//go:generate go-enum
package mapp

import (
	"fmt"
	"go/types"
	"strings"
)



// TypeFamily ENUM(basic, named, struct, pointer, slice)
type TypeFamily uint8

type TypedField interface {
	TypeFamily() TypeFamily
	Path() string
	TypeName() string
}

func extractType(s string) string {
	toReturn := s
	if strings.Contains(s, ".") {
		splitedT := strings.Split(s, ".")
		toReturn = splitedT[len(splitedT)-1]
	}

	toReturn = strings.ReplaceAll(toReturn, "*", "")
	toReturn = strings.ReplaceAll(toReturn, "[]", "")

	return toReturn
}

type Field struct {
	spec      *types.Var
	owner     *types.Struct
	fieldPath string
}

func (f Field) Name() string {
	return f.spec.Origin().Name()
}

func (f Field) FullName() string {
	return f.fieldPath + f.Name()
}

func (f Field) Fields() []Field {
	switch ft := f.spec.Type().(type) {
	case *types.Basic:
		return nil
	case *types.Named:
		_, isStruct := ft.Underlying().(*types.Struct)
		if !isStruct {
			return nil
		}

		typePath := ft.Obj().Pkg().Path()
		splitedType := strings.Split(f.spec.Origin().Type().String(), ".")
		name := splitedType[len(splitedType)-1]
		return extractFieldsFromStruct(f.FullName() + ".", typePath, name)
	case *types.Pointer:
		switch pt := ft.Elem().(type) {
		case *types.Named:
			_, isStruct := pt.Underlying().(*types.Struct)
			if !isStruct {
				return nil
			}
	
			typePath := pt.Obj().Pkg().Path()
			return extractFieldsFromStruct(f.FullName() + ".", typePath, pt.Obj().Name())
		case *types.Struct:
			extractFieldsFromStruct(f.FullName(), f.spec.Pkg().Path(), f.spec.Type().String())
		}
	}

	return nil
}


func (f Field) Type() TypedField {
	return f.resolveType(f.spec.Type())
}

func (f Field) resolveType(t types.Type) TypedField {
	switch ft := t.(type) {
	case *types.Basic:
		return BasicType{Basic: ft, f: f}
	case *types.Named:
		strt, isStruct := ft.Underlying().(*types.Struct)
		if isStruct {
			return StructType{Struct: strt, f: f}
		}
		return NamedType{Named: ft, f: f}
	case *types.Struct:
		return StructType{Struct: ft, f: f}
	case *types.Pointer:
		return PointerType{Pointer: ft, f: f}
	case *types.Slice:
		return SliceType{Slice: ft, f: f}
	default:
		panic(fmt.Sprintf("unsupported field type: %T", t))
	}
}
