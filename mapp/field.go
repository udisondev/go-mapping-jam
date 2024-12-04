//go:generate go-enum
package mapp

import (
	"fmt"
	"go/types"
	"strings"
)

type BasicType struct {
	f Field
	*types.Basic
}

type NamedType struct {
	f Field
	*types.Named
}

type StructType struct {
	f Field
	*types.Struct
}

type PointerType struct {
	f Field
	*types.Pointer
}

type SliceType struct {
	f Field
	*types.Slice
}

// TypeFamily ENUM(basic, named, struct, pointer, slice)
type TypeFamily uint8

type TypedField interface {
	TypeFamily() TypeFamily
	Path() string
	Type() string
}

func (t BasicType) TypeFamily() TypeFamily   { return FieldTypeBasic }
func (t NamedType) TypeFamily() TypeFamily   { return FieldTypeNamed }
func (t StructType) TypeFamily() TypeFamily  { return FieldTypeStruct }
func (t PointerType) TypeFamily() TypeFamily { return FieldTypePointer }
func (t SliceType) TypeFamily() TypeFamily   { return FieldTypeSlice }

func (t BasicType) Path() string {
	return t.f.spec.Pkg().Path()
}
func (t NamedType) Path() string {
	return t.f.spec.Pkg().Path()
}
func (t StructType) Path() string {
	return t.f.spec.Pkg().Path()
}
func (t PointerType) Path() string {
	return t.f.spec.Pkg().Path()
}
func (t SliceType) Path() string {
	return t.f.spec.Pkg().Path()
}

func (t BasicType) Type() string {
	return extractType(t.f.spec.Type().String())
}
func (t NamedType) Type() string {
	return extractType(t.f.spec.Type().String())
}
func (t StructType) Type() string {
	return extractType(t.f.spec.Type().String())
}
func (t PointerType) Type() string {
	return extractType(t.f.spec.Type().String())
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

func (t SliceType) Type() string {
	return extractType(t.f.spec.Type().String())
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
		return extractFieldsFromStruct(f.FullName(), typePath, name)
	case *types.Pointer:
		switch ft.Underlying().(type) {
		case *types.Struct:
			extractFieldsFromStruct(f.FullName(), f.spec.Pkg().Path(), f.spec.Type().String())
		}
	}

	return nil
}

func (p PointerType) Elem() TypedField {
	return p.f.resolveType(p.Pointer.Elem())
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
