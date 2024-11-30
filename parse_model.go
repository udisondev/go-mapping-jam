//go:generate go-enum
package main

import "fmt"

// FieldType ENUM(
// Primetive,
// Struct,
// Enum
// SliceOfStruct,
// SliceOfPrimetive,
// PointerToPrimetive
// PointerToStruct
// )
type FieldType uint8

type Mapable interface {
	fieldType() FieldType
}

func (s *Struct) fieldType() FieldType    { return FieldTypeStruct }
func (p *Primetive) fieldType() FieldType { return FieldTypePrimetive }
func (p *Slice) fieldType() FieldType {
	switch p.Of.fieldType() {
	case FieldTypePrimetive:
		return FieldTypeSliceOfPrimetive
	case FieldTypeStruct:
		return FieldTypeSliceOfStruct
	default:
		panic(fmt.Sprintf("unsupported type slice of: %T", p.Of))
	}
}
func (p *Pointer) fieldType() FieldType {
	switch p.To.fieldType() {
	case FieldTypePrimetive:
		return FieldTypePointerToPrimetive
	case FieldTypeStruct:
		return FieldTypePointerToStruct
	default:
		panic(fmt.Sprintf("unsupported type pointer to: %T", p.To))
	}
}
func (p *Enum) fieldType() FieldType { return FieldTypeEnum }

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
