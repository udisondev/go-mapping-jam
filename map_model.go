package main

import "fmt"

type Mapable interface {
	isMappable()
}

func (s *Struct) isMappable()         {}
func (p *Primetive) isMappable()      {}
func (p *StructSlice) isMappable()    {}
func (p *PrimetiveSlice) isMappable() {}

type Field struct {
	Owner *Field
	Name  string
	Desc  Mapable
}

func (f *Field) FullName() string {
	return buildFullName(f)
}

func buildFullName(f *Field) string {
	if f.Owner == nil {
		return "." + f.Name
	}
	return fmt.Sprintf("%s.%s", buildFullName(f.Owner), f.Name)
}

type Struct struct {
	Path   string
	Name   string
	Fields map[string]*Field
}

type StructSlice struct {
	Struct *Struct
}

type Primetive struct {
	Type string
}

type PrimetiveSlice struct {
	Primetive
}

type MapFunc struct {
	Name   string
	Source *Struct
	Target *Struct
	Rules map[RuleType][]Rule
}

func (s Struct) Hash() string {
	return s.Path + "." + s.Name
}

var mappersMap = make(map[string]MapFunc)
