package mapp

import (
	"go/ast"
)

type Target struct {
	spec *ast.Field
	r    Result
}

func (t Target) Name() string {
	name := t.r.Name()
	if name != "" {
		return name
	}
	return "t"
}

func (t Target) Fields() []Field {
	_, name := t.r.Type()
	return extractFieldsFromStruct(t.Name(), t.r.Path(), name)
}
