package main

// import (
// 	jen "github.com/dave/jennifer/jen"
// 	. "github.com/udisondev/go-mapping-jam/mapp"
// 	"github.com/udisondev/go-mapping-jam/rule"
// )

// type generatedMapper struct {
// 	generatedFile
// 	name      string
// 	from      Struct
// 	to        Struct
// 	isFromPrt bool
// 	isToPtr   bool
// 	rules     map[rule.Type][]rule.Any
// 	*jen.Statement
// 	fieldGenMapFuncs map[FieldType]map[FieldType]func(bl mapperBlock, target Field, source Field)
// }

// type generatedFile struct {
// 	submappers map[string]string
// }

// type mapperBlock struct {
// 	generatedMapper
// 	*jen.Group
// }

// type mappedField struct {
// 	name   string
// 	target Field
// 	source Struct
// 	mapperBlock
// }