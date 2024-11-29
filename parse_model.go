//go:generate go-enum --noprefix
package main

import "fmt"

// FieldType ENUM(
// StructType,
// PrimetiveType,
// StructSliceType,
// PrimetiveSliceType,
// PointerType,
// EnumType
// )
type FieldType uint8

type Mapable interface {
	fieldType() FieldType
}

func (s *Struct) fieldType() FieldType    { return StructType }
func (p *Primetive) fieldType() FieldType { return PrimetiveType }
func (p *Slice) fieldType() FieldType     { return StructSliceType }
func (p *Pointer) fieldType() FieldType   { return PointerType }
func (p *Enum) fieldType() FieldType      { return EnumType }

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

type Pointer struct {
	To Mapable
}

type Struct struct {
	Path   string
	Name   string
	Fields map[string]*Field
}

type Slice struct {
	Of Mapable
}

type Primetive struct {
	Type string
}

type Enum struct {
	Name string
	Path string
}

type Mapper struct {
	Name   string
	Source *Struct
	Target *Struct
	Rules  map[RuleType][]Rule
}

func (s Struct) Hash() string {
	return s.Path + "." + s.Name
}

var mappersMap = make(map[string]Mapper)
