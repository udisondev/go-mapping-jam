package mapp

import (
	"fmt"
	"go/types"
)

type StructType struct {
	f Field
	*types.Struct
}

func (t StructType) TypeFamily() TypeFamily  { return FieldTypeStruct }
func (t StructType) Path() string {
	switch st := t.f.spec.Type().(type) {
	case *types.Named:
		return st.Obj().Pkg().Path()
	case *types.Struct:
		fmt.Printf("")
	}

	panic("unsupported case")
}
func (t StructType) TypeName() string {
	return extractType(t.f.spec.Type().String())
}