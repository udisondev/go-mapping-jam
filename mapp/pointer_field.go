package mapp

import "go/types"

type PointerType struct {
	f Field
	*types.Pointer
}

func (t PointerType) TypeFamily() TypeFamily { return FieldTypePointer }
func (t PointerType) TypeName() string {
	return extractType(t.f.spec.Type().String())
}
func (p PointerType) Elem() TypedField {
	return p.f.resolveType(p.Pointer.Elem())
}
func (t PointerType) Path() string {
	return t.f.spec.Pkg().Path()
}
