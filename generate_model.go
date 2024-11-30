package main

import jen "github.com/dave/jennifer/jen"

type generatedMapper struct {
	generatedFile
	name      string
	from      Struct
	to        Struct
	isFromPrt bool
	isToPtr   bool
	rules     map[RuleType][]Rule
	*jen.Statement
	fieldGenMapFuncs map[FieldType]map[FieldType]func(bl mapperBlock, target Field, source Field)
}

type generatedFile struct {
	submappers map[string]string
}

type mapperBlock struct {
	generatedMapper
	*jen.Group
}

type mappedField struct {
	name   string
	target Field
	source Struct
	mapperBlock
}
