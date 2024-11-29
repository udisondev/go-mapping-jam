package main

import jen "github.com/dave/jennifer/jen"

type generatedMapper struct {
	name       string
	from       *Struct
	to         *Struct
	isFromPrt  bool
	isToPtr    bool
	file       *jen.File
	rules      map[RuleType][]Rule
	statement  *jen.Statement
	submappers map[string]string
}

type mapperBlock struct {
	from  func() *Struct
	to    func() *Struct
	group *jen.Group
}

type mappedField struct {
	name       string
	file       func() *jen.File
	field      *Field
	source     func() *Struct
	group      func() *jen.Group
	rules      func() map[RuleType][]Rule
	submappers func() map[string]string
}

type mappingCase string
