package mapp

import (
	"go/ast"
	"strings"
)

type Comment struct {
	spec *ast.Comment
}

func (c Comment) Value() string {
	val, _ := strings.CutPrefix(c.spec.Text, "//")
	return val
}