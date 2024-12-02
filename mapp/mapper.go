package mapp

import (
	"go/ast"
	"regexp"
)

var ruleReg = regexp.MustCompile(`^@(qual|enum)=`)

type Mapper struct {
	spec    *ast.Field
	imports []Import
}

func (m *Mapper) Name() string {
	return m.spec.Names[0].Name
}

func (m *Mapper) Comments() []Comment {
	comments := make([]Comment, 0, len(m.spec.Doc.List))
	for _, v := range m.spec.Doc.List {
		comments = append(comments, Comment{spec: v})
	}
	return comments
}

func (m *Mapper) Rules() []Rule {
	rules := []Rule{}
	for _, c := range m.Comments() {
		val := c.Value()
		if ruleReg.MatchString(val) {
			rules = append(rules, Rule{
				spec: val,
			})
		}
	}

	return rules
}

func (m *Mapper) Params() []Param {
	fnT, ok := m.spec.Type.(*ast.FuncType)
	if !ok {
		panic("is not a ast.FuncType")
	}

	params := make([]Param, 0, len(fnT.Params.List))
	for _, p := range fnT.Params.List {
		params = append(params, Param{spec: p, imports: m.imports})
	}

	return params
}

func (m *Mapper) Results() []Result {
	fnT, ok := m.spec.Type.(*ast.FuncType)
	if !ok {
		panic("is not a ast.FuncType")
	}

	params := make([]Result, 0, len(fnT.Results.List))
	for _, p := range fnT.Results.List {
		params = append(params, Result{spec: p, imports: m.imports})
	}

	return params
}

func (m *Mapper) Source() Source {
	return Source{
		spec: m.Params()[0].spec,
		Param: m.Params()[0],
	}
}

func (m *Mapper) Target() Target {
	return Target{
		spec: m.Results()[0].spec,
		Result: m.Results()[0],
	}
}
