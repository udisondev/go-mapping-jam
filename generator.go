package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	jen "github.com/dave/jennifer/jen"
)

type generatedMapper struct {
	name       string
	from       *Struct
	to         *Struct
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
	from       func() *Struct
	to         func() *Struct
	group      func() *jen.Group
	rules      func() map[RuleType][]Rule
	submappers func() map[string]string
}

type mappingCase string

const (
	TargetPrimetive_SourcePrimetive       mappingCase = "TargetPrimetiveType_SourcePrimetiveType"
	TargetPrimetive_SourcePtrPrimetive    mappingCase = "TargetPrimetiveType_SourcePtrPrimetiveType"
	TargetPtrPrimetive_SourcePrimetive    mappingCase = "TargetPtrPrimetiveType_SourcePrimetiveType"
	TargetPtrPrimetive_SourcePtrPrimetive mappingCase = "TargetPtrPrimetiveType_SourcePtrPrimetiveType"
	TargetStruct_SourceStruct             mappingCase = "TargetStructType_SourceStructType"
	TargetStruct_SourcePtrStruct          mappingCase = "TargetStructType_SourcePtrStructType"
	TargetPtrStruct_SourceStruct          mappingCase = "TargetPtrStructType_SourceStructType"
	TargetPtrStruct_SourcePtrStruct       mappingCase = "TargetPtrStructType_SourcePtrStructType"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

func (gm *generatedMapper) generateMapFunc() {
	gm.initSignature()
	gm.initBody()
}

func (bl *mapperBlock) initTarget() {
	if bl.to().Path == currentPath {
		bl.group.Id("target").Op(":=").Id(bl.to().Name + "{}")
	} else {
		bl.group.Id("target").Op(":=").Qual(bl.to().Path, bl.to().Name+"{}")
	}
}

func (mf *mappedField) mapField() {
	sourceFieldName := mf.name
	if qr, ok := mf.findQualRule(); ok && qr.SourceName != "" && qr.MName == "" {
		sourceFieldName = qr.SourceName
	} else if ok && qr.MName != "" && qr.MPath != "" {
		if qr.SourceName != "" {
			sourceFieldName = qr.SourceName
		}
		mf.group().Id("target").
			Dot(mf.name).
			Op("=").
			Qual(qr.MPath, qr.MName).
			Call(jen.Id("src").Dot(sourceFieldName))
		return
	} else if ok && qr.MName != "" {
		if qr.SourceName != "" {
			sourceFieldName = qr.SourceName
		}
		mf.group().Id("target").
			Dot(mf.name).
			Op("=").
			Id(qr.MName).
			Call(jen.Id("src").Dot(sourceFieldName))
		return
	}

	sourceField, ok := mf.from().Fields[sourceFieldName]
	if !ok {
		log.Fatalf("source field not found for target: %s", mf.field.FullName())
	}

	switch mf.resolveFieldMapper(sourceField) {
	case TargetPrimetive_SourcePrimetive:
		mf.group().Id("target").Dot(mf.name).Op("=").Id("src").Dot(sourceFieldName)
	case TargetPrimetive_SourcePtrPrimetive:
		mf.genPrimetivePtrPrimetive(sourceFieldName)
	case TargetPtrPrimetive_SourcePrimetive:
		mf.genPtrPrimetivePrimetive(sourceFieldName)
	case TargetPtrPrimetive_SourcePtrPrimetive:
	case TargetStruct_SourceStruct:
		mf.genStructStructMapping(sourceFieldName, sourceField)
	case TargetStruct_SourcePtrStruct:
	case TargetPtrStruct_SourceStruct:
	case TargetPtrStruct_SourcePtrStruct:
	}
}

func (mf *mappedField) genPrimetivePtrPrimetive(sourceFieldName string) {
	mf.group().
		If(
			jen.Id("src").Dot(sourceFieldName).Op("!=").Nil(),
		).
		Block(
			jen.Id("target").Dot(mf.field.Name).Op("=").Add(jen.Op("*")).Id("src").Dot(sourceFieldName),
		)
}

func (mf *mappedField) genPtrPrimetivePrimetive(sourceFieldName string) {
	mf.group().Id("target").Dot(mf.name).Op("=").Add(jen.Op("&")).Id("src").Dot(sourceFieldName)
}

func (mf *mappedField) genStructStructMapping(sourceFieldName string, sourceField *Field) {
	nestedSourceStruct, _ := sourceField.Desc.(*Struct)
	hash := nestedSourceStruct.Hash() + mf.field.Desc.(*Struct).Hash()
	methodName, ok := mf.submappers()[hash]
	if !ok {
		methodName = genRandomName(15)
		mf.submappers()[hash] = methodName
	}
	mf.group().Id("target").Dot(mf.name).Op("=").Id(methodName).Call(jen.Id("src").Dot(sourceFieldName))

	if !ok {
		sbm := generatedMapper{
			name:       methodName,
			from:       sourceField.Desc.(*Struct),
			to:         mf.field.Desc.(*Struct),
			file:       mf.file(),
			rules:      mf.rules(),
			submappers: mf.submappers(),
		}
		sbm.generateMapFunc()
	}
}

func (mf *mappedField) findQualRule() (QualRule, bool) {
	for _, v := range mf.rules()[Qual] {
		qr, ok := v.(QualRule)
		if !ok {
			panic("is not qual rule")
		}

		if qr.TargetName == mf.field.FullName() {
			return qr, true
		}
	}

	return QualRule{}, false
}

func (gm *generatedMapper) initBody() {
	gm.statement.BlockFunc(func(gr *jen.Group) {
		gbl := mapperBlock{
			from:  func() *Struct { return gm.from },
			to:    func() *Struct { return gm.to },
			group: gr,
		}
		gbl.initTarget()
		for n, f := range gm.to.Fields {
			field := mappedField{
				name:       n,
				field:      f,
				file:       func() *jen.File { return gm.file },
				from:       func() *Struct { return gm.from },
				to:         func() *Struct { return gm.to },
				group:      func() *jen.Group { return gr },
				rules:      func() map[RuleType][]Rule { return gm.rules },
				submappers: func() map[string]string { return gm.submappers },
			}
			field.mapField()
		}
		gr.Return(jen.Id("target"))
	})
}

func generateCodeWithJennifer(outputFile string, mapFuncs map[string]Mapper) {
	f := jen.NewFile("mapper")
	f.Comment("Code generated by go-mapping-jam. DO NOT EDIT.")

	submappers := make(map[string]string)
	for _, mapFunc := range mapFuncs {
		gm := generatedMapper{
			name:       mapFunc.Name,
			from:       mapFunc.Source,
			to:         mapFunc.Target,
			file:       f,
			rules:      mapFunc.Rules,
			submappers: submappers,
		}
		gm.generateMapFunc()
	}

	err := f.Save(outputFile)
	if err != nil {
		log.Fatalf("failed to save file: %v", err)
	}
}

func (gf *mappedField) resolveFieldMapper(sourceField *Field) mappingCase {
	buildMappingCase := func(targetType, sourceType FieldType, isTargetPtr, isSourcePtr bool) mappingCase {
		ptrStr := func(isPrt bool) string {
			if isPrt {
				return "Ptr"
			}
			return ""
		}

		return mappingCase(fmt.Sprintf("Target%s%s_Source%s%s", ptrStr(isTargetPtr), targetType, ptrStr(isSourcePtr), sourceType))
	}

	if gf.field.Desc.fieldType() != PointerType && sourceField.Desc.fieldType() != PointerType {
		return buildMappingCase(gf.field.Desc.fieldType(), sourceField.Desc.fieldType(), false, false)
	}

	if gf.field.Desc.fieldType() == PointerType && sourceField.Desc.fieldType() != PointerType {
		tarPtr, ok := gf.field.Desc.(*Pointer)
		if !ok {
			panic("is not a pointer")
		}
		return buildMappingCase(tarPtr.Ref.fieldType(), sourceField.Desc.fieldType(), true, false)
	}

	if gf.field.Desc.fieldType() == PointerType && sourceField.Desc.fieldType() == PointerType {
		tarPtr, ok := gf.field.Desc.(*Pointer)
		if !ok {
			panic("is not a pointer")
		}

		srcPtr, ok := sourceField.Desc.(*Pointer)
		if !ok {
			panic("is not a pointer")
		}
		return buildMappingCase(tarPtr.Ref.fieldType(), srcPtr.Ref.fieldType(), true, true)
	}

	if gf.field.Desc.fieldType() != PointerType && sourceField.Desc.fieldType() == PointerType {
		srcPtr, ok := sourceField.Desc.(*Pointer)
		if !ok {
			panic("is not a pointer")
		}
		return buildMappingCase(gf.field.Desc.fieldType(), srcPtr.Ref.fieldType(), false, true)
	}

	return ""
}

func (gm *generatedMapper) initSignature() {
	gm.statement = gm.file.Func().Id(gm.name)
	if gm.from.Path == currentPath {
		gm.statement.Params(jen.Id("src").Id(gm.from.Name))
	} else {
		gm.statement.Params(jen.Id("src").Qual(gm.from.Path, gm.from.Name))
	}

	if gm.to.Path == currentPath {
		gm.statement.Id(gm.to.Name)
	} else {
		gm.statement.Qual(gm.to.Path, gm.to.Name)
	}
}

func genRandomName(length int) string {
	seed := time.Now().UnixNano()

	src := rand.NewPCG(uint64(seed), uint64(seed>>32))
	r := rand.New(src)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.IntN(len(charset))]
	}
	return string(result)
}
