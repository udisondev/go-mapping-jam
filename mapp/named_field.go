package mapp

import "go/types"

type NamedType struct {
	f Field
	*types.Named
}

func (t NamedType) TypeFamily() TypeFamily   { return FieldTypeNamed }
func (t NamedType) Path() string {
	return t.f.spec.Pkg().Path()
}

func (t NamedType) TypeName() string {
	return extractType(t.f.spec.Type().String())
}