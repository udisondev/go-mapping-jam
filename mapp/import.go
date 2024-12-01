package mapp

import (
	"go/ast"
	"path/filepath"
	"strings"
)

type Import struct {
	spec *ast.ImportSpec
}

func (i *Import) Alias() string {
	if i.spec.Name != nil {
		return i.spec.Name.Name
	}

	base := filepath.Base(strings.ReplaceAll(i.spec.Path.Value, "\"", ""))
	return base
}

func (i *Import) Path() string {
	return i.spec.Path.Value
}
