package mapp

import (
	"fmt"
	"go/types"
)

type SliceType struct {
	f Field
	*types.Slice
}

func (s SliceType) TypeFamily() TypeFamily   { return FieldTypeSlice }
func (s SliceType) Path() string {
	sl, ok := s.f.spec.Type().(*types.Slice)
	if !ok {
		panic("is not a slice")
	}
	switch st := sl.Elem().(type) {
	case *types.Basic:
		return ""
	case *types.Named:
		return st.Obj().Pkg().Path()
	case *types.Struct:
		fmt.Printf("")
	}

	panic("unsupported case")
}
func (s SliceType) TypeName() string {
	return extractType(s.f.spec.Type().String())
}

func (s SliceType) Elem() TypedField {
	return s.f.resolveType(s.Slice.Elem())
}