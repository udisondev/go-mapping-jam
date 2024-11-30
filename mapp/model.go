//go:generate go-enum
package mapp

import (
	"fmt"

	"github.com/udisondev/go-mapping-jam/rule"
)

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
	FieldType() FieldType
}

func (s Struct) FieldType() FieldType    { return FieldTypeStruct }
func (p Primetive) FieldType() FieldType { return FieldTypePrimetive }
func (p Slice) FieldType() FieldType {
	switch p.Of.FieldType() {
	case FieldTypePrimetive:
		return FieldTypeSliceOfPrimetive
	case FieldTypeStruct:
		return FieldTypeSliceOfStruct
	default:
		panic(fmt.Sprintf("unsupported type slice of: %T", p.Of))
	}
}
func (p Pointer) FieldType() FieldType {
	switch p.To.FieldType() {
	case FieldTypePrimetive:
		return FieldTypePointerToPrimetive
	case FieldTypeStruct:
		return FieldTypePointerToStruct
	default:
		panic(fmt.Sprintf("unsupported type pointer to: %T", p.To))
	}
}
func (p Enum) FieldType() FieldType { return FieldTypeEnum }

type Field struct {
	Owner *Field
	Name  string
	Desc  Mapable
}

func (f Field) FullName() string {
	return buildFullName(f)
}

func buildFullName(f Field) string {
	if f.Owner == nil {
		return "." + f.Name
	}
	return fmt.Sprintf("%s.%s", buildFullName(*f.Owner), f.Name)
}

type Pointer struct {
	To Mapable
}

type Struct struct {
	Path   string
	Name   string
	Fields map[string]Field
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
	Source Struct
	Target Struct
	Rules  map[rule.Type][]rule.Any
}

func (s Struct) Hash() string {
	return s.Path + "." + s.Name
}

