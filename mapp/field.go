//go:generate go-enum
package mapp

import (
	"fmt"
	"go/types"
	"strings"
)

type BasicType struct {
	*types.Basic
}

type NamedType struct {
	*types.Named
}

type StructType struct {
	*types.Struct
}

type PointerType struct {
	*types.Pointer
}

type SliceType struct {
	*types.Slice
}

// FieldType ENUM(basic, named, struct, pointer, slice)
type FieldType uint8

type TypedField interface {
	ShortType() FieldType
}

func (t BasicType) ShortType() FieldType   { return FieldTypeBasic }
func (t NamedType) ShortType() FieldType   { return FieldTypeNamed }
func (t StructType) ShortType() FieldType  { return FieldTypeStruct }
func (t PointerType) ShortType() FieldType { return FieldTypePointer }
func (t SliceType) ShortType() FieldType   { return FieldTypeSlice }

type Field struct {
	spec      *types.Var
	owner     *types.Struct
	fieldPath string
}

func (f *Field) Name() string {
	return f.spec.Origin().Name()
}

func (f *Field) FullName() string {
	return f.fieldPath + "." + f.Name()
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
		return extractFieldsFromStruct(f.FullName(), path, name)
	case *types.Pointer:
		switch ft.Underlying().(type) {
		case *types.Struct:
			extractFieldsFromStruct(f.FullName(), f.spec.Pkg().Path(), f.spec.Type().String())
		}
	}

	return nil
}

func (f *Field) Type() TypedField {
	switch ft := f.spec.Type().(type) {
	case *types.Basic:
		return BasicType{Basic: ft}
	case *types.Named:
		strt, isStruct := ft.Underlying().(*types.Struct)
		if isStruct {
			return StructType{Struct: strt}
		}
		return NamedType{Named: ft}
	case *types.Struct:
		return StructType{Struct: ft}
	case *types.Pointer:
		return PointerType{Pointer: ft}
	case *types.Slice:
		return SliceType{Slice: ft}
	default:
		panic(fmt.Sprintf("unsupported field type: %T", f.spec.Type()))
	}
}
