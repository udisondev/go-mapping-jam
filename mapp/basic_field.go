package mapp

import "go/types"

type BasicType struct {
	f Field
	*types.Basic
}
func (t BasicType) TypeFamily() TypeFamily   { return FieldTypeBasic }
func (t BasicType) Path() string {
	return t.f.spec.Pkg().Path()
}
func (t BasicType) TypeName() string {
	return extractType(t.f.spec.Type().String())
}
