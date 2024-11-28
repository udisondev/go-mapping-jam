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

const (
	TargetPrimetive_SourcePrimetive                 mappingCase = "TargetPrimetiveType_SourcePrimetiveType"
	TargetPrimetive_SourcePtrPrimetive              mappingCase = "TargetPrimetiveType_SourcePtrPrimetiveType"
	TargetPtrPrimetive_SourcePrimetive              mappingCase = "TargetPtrPrimetiveType_SourcePrimetiveType"
	TargetPtrPrimetive_SourcePtrPrimetive           mappingCase = "TargetPtrPrimetiveType_SourcePtrPrimetiveType"
	TargetPrimetiveSlice_SourcePrimetiveSlice       mappingCase = "TargetPrimetiveSliceType_SourcePrimetiveSliceType"
	TargetPrimetiveSlice_SourcePtrPrimetiveSlice    mappingCase = "TargetPrimetiveSliceType_SourcePtrPrimetiveSliceType"
	TargetPtrPrimetiveSlice_SourcePrimetiveSlice    mappingCase = "TargetPtrPrimetiveSliceType_SourcePrimetiveSliceType"
	TargetPtrPrimetiveSlice_SourcePtrPrimetiveSlice mappingCase = "TargetPtrPrimetiveSliceType_SourcePtrPrimetiveSliceType"
	TargetStruct_SourceStruct                       mappingCase = "TargetStructType_SourceStructType"
	TargetStruct_SourcePtrStruct                    mappingCase = "TargetStructType_SourcePtrStructType"
	TargetPtrStruct_SourceStruct                    mappingCase = "TargetPtrStructType_SourceStructType"
	TargetPtrStruct_SourcePtrStruct                 mappingCase = "TargetPtrStructType_SourcePtrStructType"
)