//go:generate go-enum
package main

// import (
// 	"fmt"
// 	"go/types"
// 	"log"
// 	"strings"

// 	"github.com/udisondev/go-mapping-jam/rule"
// 	"golang.org/x/tools/go/packages"
// )

// // FieldType ENUM(
// // Primetive,
// // Struct,
// // Enum
// // SliceOfStruct,
// // SliceOfPrimetive,
// // PointerToPrimetive
// // PointerToStruct
// // )
// type FieldType uint8

// type Mapable interface {
// 	FieldType() FieldType
// }

// func (s Struct) FieldType() FieldType    { return FieldTypeStruct }
// func (p Primetive) FieldType() FieldType { return FieldTypePrimetive }
// func (p Slice) FieldType() FieldType {
// 	switch p.Of.FieldType() {
// 	case FieldTypePrimetive:
// 		return FieldTypeSliceOfPrimetive
// 	case FieldTypeStruct:
// 		return FieldTypeSliceOfStruct
// 	default:
// 		panic(fmt.Sprintf("unsupported type slice of: %T", p.Of))
// 	}
// }
// func (p Pointer) FieldType() FieldType {
// 	switch p.To.FieldType() {
// 	case FieldTypePrimetive:
// 		return FieldTypePointerToPrimetive
// 	case FieldTypeStruct:
// 		return FieldTypePointerToStruct
// 	default:
// 		panic(fmt.Sprintf("unsupported type pointer to: %T", p.To))
// 	}
// }
// func (p Enum) FieldType() FieldType { return FieldTypeEnum }

// type Field struct {
// 	Owner *Field
// 	Name  string
// 	Desc  Mapable
// }

// func (f Field) FullName() string {
// 	return buildFullName(f)
// }

// func buildFullName(f Field) string {
// 	if f.Owner == nil {
// 		return "." + f.Name
// 	}
// 	return fmt.Sprintf("%s.%s", buildFullName(*f.Owner), f.Name)
// }

// type Pointer struct {
// 	To Mapable
// }

// type Struct struct {
// 	Path   string
// 	Name   string
// 	Fields map[string]Field
// }

// type Slice struct {
// 	Of Mapable
// }

// type Primetive struct {
// 	Type string
// }

// type Enum struct {
// 	Name string
// 	Path string
// }

// type Mapper struct {
// 	Name   string
// 	Source Struct
// 	Target Struct
// 	Rules  map[rule.Type][]rule.Any
// }

// func (s Struct) Hash() string {
// 	return s.Path + "." + s.Name
// }

// func (m *Mapper) InitRoot(pkgs map[string]*packages.Package, str Struct, packFunc func(dir string) *packages.Package) {
// 	pkg, ok := pkgs[str.Path]
// 	if !ok {
// 		pkg = packFunc(dirByPath(str.Path))
// 		pkgs[str.Path] = pkg
// 	}

// 	obj := pkg.Types.Scope().Lookup(str.Name)
// 	if obj == nil {
// 		log.Fatalf("type %s not found in package %s", str.Name, str.Path)
// 	}

// 	namedType, ok := obj.Type().(*types.Named)
// 	if !ok {
// 		log.Fatalf("%s is not a named type", str.Name)
// 	}

// 	structType, ok := namedType.Underlying().(*types.Struct)
// 	if !ok {
// 		log.Fatalf("%s is not a struct", str.Name)
// 	}

// 	for i := 0; i < structType.NumFields(); i++ {
// 		field := structType.Field(i)
// 		fieldName := field.Name()
// 		str.Fields[fieldName] = m.buildField(nil, fieldName, field)
// 	}
// }

// func (m *Mapper) buildField(owner *Field, fieldName string, field *types.Var) Field {
// 	switch t := field.Type().(type) {
// 	case *types.Basic:
// 		return Field{
// 			Owner: owner,
// 			Name:  fieldName,
// 			Desc: &Primetive{
// 				Type: t.Name(),
// 			}}
// 	case *types.Named:
// 		fmt.Printf("Name: %s\n", t.Obj().Name())
// 		fmt.Printf("Pkg name: %s\n", t.Obj().Pkg().Name())
// 		fmt.Printf("Pkg path: %s\n", t.Obj().Pkg().Path())
// 		structType, ok := t.Underlying().(*types.Struct)
// 		if !ok {
// 			return Field{
// 				Owner: owner,
// 				Name:  fieldName,
// 				Desc: &Enum{
// 					Name: t.Obj().Name(),
// 					Path: t.Obj().Pkg().Name(),
// 				},
// 			}
// 		}

// 		fs := Field{
// 			Owner: owner,
// 			Name:  fieldName,
// 			Desc: &Struct{
// 				Path:   t.Obj().Pkg().Path(),
// 				Name:   t.Obj().Name(),
// 				Fields: make(map[string]Field),
// 			},
// 		}

// 		for i := 0; i < structType.NumFields(); i++ {
// 			subField := structType.Field(i)
// 			subFieldName := subField.Name()
// 			fs.Desc.(*Struct).Fields[subFieldName] = m.buildField(&fs, subFieldName, subField)
// 		}

// 		return fs

// 	case *types.Slice:
// 		switch slt := t.Elem().(type) {
// 		case *types.Basic:
// 			return Field{
// 				Owner: owner,
// 				Name:  fieldName,
// 				Desc: &Slice{Of: &Primetive{
// 					Type: slt.Name(),
// 				}}}

// 		case *types.Named:
// 			structType, ok := slt.Underlying().(*types.Struct)
// 			if !ok {
// 				return Field{
// 					Owner: owner,
// 					Name:  fieldName,
// 					Desc: &Slice{Of: &Primetive{
// 						Type: slt.Obj().Name(),
// 					}}}
// 			}

// 			fs := Field{
// 				Owner: owner,
// 				Name:  fieldName,
// 				Desc: &Slice{
// 					Of: &Struct{
// 						Path:   slt.Obj().Pkg().Path(),
// 						Name:   slt.Obj().Name(),
// 						Fields: make(map[string]Field),
// 					}}}

// 			for i := 0; i < structType.NumFields(); i++ {
// 				subField := structType.Field(i)
// 				subFieldName := subField.Name()
// 				fs.Desc.(*Slice).Of.(*Struct).Fields[subFieldName] = m.buildField(&fs, subFieldName, subField)
// 			}

// 			return fs
// 		}
// 	case *types.Pointer:
// 		switch pt := t.Elem().(type) {
// 		case *types.Basic:
// 			return Field{
// 				Owner: owner,
// 				Name:  fieldName,
// 				Desc: &Pointer{
// 					To: &Primetive{
// 						Type: pt.Name(),
// 					},
// 				},
// 			}
// 		case *types.Named:
// 			structType, ok := pt.Underlying().(*types.Struct)
// 			if !ok {
// 				return Field{
// 					Owner: owner,
// 					Name:  fieldName,
// 					Desc: &Pointer{
// 						To: &Primetive{
// 							Type: pt.Obj().Name(),
// 						},
// 					},
// 				}
// 			}

// 			fs := Field{
// 				Owner: owner,
// 				Name:  fieldName,
// 				Desc: &Pointer{
// 					To: &Struct{
// 						Path:   pt.Obj().Pkg().Path(),
// 						Name:   pt.Obj().Name(),
// 						Fields: make(map[string]Field),
// 					}},
// 			}

// 			for i := 0; i < structType.NumFields(); i++ {
// 				subField := structType.Field(i)
// 				subFieldName := subField.Name()
// 				fs.Desc.(*Pointer).To.(*Struct).Fields[subFieldName] = m.buildField(&fs, subFieldName, subField)
// 			}

// 			return fs

// 		case *types.Slice:

// 			switch slt := pt.Elem().Underlying().(type) {
// 			case *types.Basic:
// 				return Field{
// 					Owner: owner,
// 					Name:  fieldName,
// 					Desc: &Slice{Of: &Primetive{
// 						Type: slt.Name(),
// 					}}}
// 			case *types.Named:
// 				structType, ok := slt.Underlying().(*types.Struct)
// 				if !ok {
// 					return Field{
// 						Owner: owner,
// 						Name:  fieldName,
// 						Desc: &Slice{Of: &Primetive{
// 							Type: slt.Obj().Name(),
// 						}}}
// 				}

// 				fs := Field{
// 					Owner: owner,
// 					Name:  fieldName,
// 					Desc: &Slice{
// 						Of: &Struct{
// 							Path:   slt.Obj().Pkg().Path(),
// 							Name:   slt.Obj().Name(),
// 							Fields: make(map[string]Field),
// 						}}}

// 				for i := 0; i < structType.NumFields(); i++ {
// 					subField := structType.Field(i)
// 					subFieldName := subField.Name()
// 					fs.Desc.(*Struct).Fields[subFieldName] = m.buildField(&fs, subFieldName, subField)
// 				}

// 				return fs
// 			}
// 		}
// 	}

// 	panic(fmt.Sprintf("unknown field type: %v", field.Type()))
// }

// func dirByPath(p string) string {
// 	return strings.ReplaceAll(p, "github.com/udisondev/go-mapping-jam", "./")
// }
