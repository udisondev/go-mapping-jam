package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	jen "github.com/dave/jennifer/jen"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

func (gm *generatedMapper) generateMapFunc() {
	gm.initSignature()
	gm.initBody()
}

func (bl *mapperBlock) initTarget() {
	if bl.to().Path == currentPath {
		bl.Id("target").Op(":=").Id(bl.to().Name + "{}")
	} else {
		bl.Id("target").Op(":=").Qual(bl.to().Path, bl.to().Name+"{}")
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
		mf.Id("target").
			Dot(mf.name).
			Op("=").
			Qual(qr.MPath, qr.MName).
			Call(jen.Id("src").Dot(sourceFieldName))
		return
	} else if ok && qr.MName != "" {
		if qr.SourceName != "" {
			sourceFieldName = qr.SourceName
		}
		mf.Id("target").
			Dot(mf.name).
			Op("=").
			Id(qr.MName).
			Call(jen.Id("src").Dot(sourceFieldName))
		return
	}

	sourceField, ok := mf.source().Fields[sourceFieldName]
	if !ok {
		log.Fatalf("source field not found for target: %s", mf.target.FullName())
	}

	mappingGenFunc, ok := mf.fieldGenMapFuncs[mf.target.Desc.fieldType()][sourceField.Desc.fieldType()]
	if !ok {
		panic(fmt.Sprintf("unable to map from '%s' to '%s'", sourceField.Desc.fieldType(), mf.target.Desc.fieldType()))
	}

	mappingGenFunc(mf.mapperBlock, mf.target, sourceField)
}

func genPrimetivePtrPrimetive(bl *mapperBlock, target, source *Field) {
	bl.
		If(
			jen.Id("src").Dot(source.Name).Op("!=").Nil(),
		).
		Block(
			jen.Id("target").Dot(target.Name).Op("=").Add(jen.Op("*")).Id("src").Dot(source.Name),
		)
}

func genPrimetivePrimetive(bl *mapperBlock, target, source *Field) {
	bl.Id("target").Dot(target.Name).Op("=").Id("src").Dot(source.Name)
}

func genPtrPrimetivePrimetive(bl *mapperBlock, target, source *Field) {
	bl.Id("target").Dot(target.Name).Op("=").Add(jen.Op("&")).Id("src").Dot(source.Name)
}

func genPtrStructPtrStructMapping(bl *mapperBlock, target, source *Field) {
	nestedSourceStruct, ok := source.Desc.(*Pointer).To.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	targetField, ok := target.Desc.(*Pointer).To.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	hash := nestedSourceStruct.Hash() + targetField.Hash()
	methodName, ok := bl.submappers[hash]
	if !ok {
		methodName = genRandomName(15)
		bl.submappers[hash] = methodName
	}

	bl.
		If(
			jen.Id("src").Dot(source.Name).Op("!=").Nil(),
		).
		Block(
			jen.Id(methodName+"Result").Op(":=").Id(methodName).Call(jen.Add(jen.Op("*").Id("src").Dot(source.Name))),
			jen.Id("target").Dot(target.Name).Op("=").Add(jen.Op("&")).Id(methodName+"Result"),
		)

	if !ok {
		sbm := generatedMapper{
			generatedFile: bl.generatedFile,
			name:          methodName,
			from:          nestedSourceStruct,
			to:            targetField,
			isFromPrt:     true,
			isToPtr:       true,
			rules:         bl.rules,
		}
		sbm.generateMapFunc()
	}
}

func genPtrStructStructMapping(bl *mapperBlock, target, source *Field) {
	nestedSourceStruct, ok := source.Desc.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	targetField, ok := target.Desc.(*Pointer).To.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	hash := nestedSourceStruct.Hash() + targetField.Hash()
	methodName, ok := bl.submappers[hash]
	if !ok {
		methodName = genRandomName(15)
		bl.submappers[hash] = methodName
	}

	bl.Id(methodName + "Result").Op(":=").Id(methodName).Call(jen.Id("src").Dot(source.Name))
	bl.Id("target").Dot(target.Name).Op("=").Add(jen.Op("&")).Id(methodName + "Result")

	if !ok {
		sbm := generatedMapper{
			generatedFile: bl.generatedFile,
			name:          methodName,
			from:          nestedSourceStruct,
			to:            targetField,
			isFromPrt:     true,
			rules:         bl.rules,
		}
		sbm.generateMapFunc()
	}
}

func genStructPtrStructMapping(bl *mapperBlock, target, source *Field) {
	nestedSourceStruct, ok := source.Desc.(*Pointer).To.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	targetField, ok := target.Desc.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	hash := nestedSourceStruct.Hash() + targetField.Hash()
	methodName, ok := bl.submappers[hash]
	if !ok {
		methodName = genRandomName(15)
		bl.submappers[hash] = methodName
	}

	bl.
		If(
			jen.Id("src").Dot(source.Name).Op("!=").Nil(),
		).
		Block(
			jen.Id("target").Dot(target.Name).Op("=").Id(methodName).Call(jen.Add(jen.Op("*").Id("src").Dot(source.Name))),
		)

	if !ok {
		sbm := generatedMapper{
			generatedFile: bl.generatedFile,
			name:          methodName,
			from:          nestedSourceStruct,
			to:            targetField,
			isFromPrt:     true,
			rules:         bl.rules,
		}
		sbm.generateMapFunc()
	}
}

func genStructStructMapping(bl *mapperBlock, target, source *Field) {
	nestedSourceStruct, _ := source.Desc.(*Struct)
	hash := nestedSourceStruct.Hash() + target.Desc.(*Struct).Hash()
	methodName, ok := bl.submappers[hash]
	if !ok {
		methodName = genRandomName(15)
		bl.submappers[hash] = methodName
	}
	bl.Id("target").Dot(target.Name).Op("=").Id(methodName).Call(jen.Id("src").Dot(source.Name))

	if !ok {
		sbm := generatedMapper{
			generatedFile: bl.generatedFile,
			name:          methodName,
			from:          source.Desc.(*Struct),
			to:            target.Desc.(*Struct),
			rules:         bl.rules,
		}
		sbm.generateMapFunc()
	}
}

func genStructSliceStructSliceMapping(bl *mapperBlock, target, source *Field) {
	targetStruct := target.Desc.(*Slice).Of.(*Struct)

	nestedSourceStruct := source.Desc.(*Slice).Of.(*Struct)
	hash := nestedSourceStruct.Hash() + targetStruct.Hash()
	methodName, ok := bl.submappers[hash]
	if !ok {
		methodName = genRandomName(15)
		bl.submappers[hash] = methodName
	}
	bl.Id("target"+target.Name+"Slice").Op(":=").Make(jen.Index().Qual(targetStruct.Path, targetStruct.Name), jen.Lit(0), jen.Len(jen.Id("target").Dot(target.Name)))
	bl.
		For(
			jen.List(jen.Id("_"), jen.Id("it")).Op(":=").Range().Id("src").Dot(source.Name),
		).
		Block(
			jen.Id("target"+target.Name+"Slice").Op("=").Append(jen.Id("target"+target.Name+"Slice"), jen.Id(methodName).Call(jen.Id("it"))),
		)

	bl.Id("target").Dot(target.Name).Op("=").Id("target" + target.Name + "Slice")

	if !ok {
		sbm := generatedMapper{
			generatedFile: bl.generatedFile,
			name:          methodName,
			from:          nestedSourceStruct,
			to:            targetStruct,
			rules:         bl.rules,
		}
		sbm.generateMapFunc()
	}
}

func (mf *mappedField) findQualRule() (QualRule, bool) {
	for _, v := range mf.rules[RuleTypeQual] {
		qr, ok := v.(QualRule)
		if !ok {
			panic("is not qual rule")
		}

		if qr.TargetName == mf.target.FullName() {
			return qr, true
		}
	}

	return QualRule{}, false
}

func (gm *generatedMapper) initBody() {
	gm.BlockFunc(func(gr *jen.Group) {
		gbl := mapperBlock{
			from:            func() *Struct { return gm.from },
			to:              func() *Struct { return gm.to },
			generatedMapper: gm,
			Group:           gr,
		}
		gbl.initTarget()
		for n, f := range gm.to.Fields {
			field := mappedField{
				name:        n,
				target:      f,
				source:      func() *Struct { return gm.from },
				mapperBlock: &gbl,
			}
			field.mapField()
		}
		gr.Return(jen.Id("target"))
	})
}

func generateCodeWithJennifer(outputFile string, mapFuncs map[string]Mapper) {
	f := jen.NewFile("mapper")
	g := generatedFile{
		submappers: make(map[string]string),
	}

	fieldGenMap := map[FieldType]map[FieldType]func(bl *mapperBlock, target *Field, source *Field){
		FieldTypePrimetive: {
			FieldTypePrimetive:          genPrimetivePrimetive,
			FieldTypePointerToPrimetive: genPrimetivePtrPrimetive,
		},
		FieldTypeStruct: {
			FieldTypeStruct:          genStructStructMapping,
			FieldTypePointerToStruct: genStructPtrStructMapping,
		},
		FieldTypeEnum: {
			FieldTypeEnum: genPrimetivePrimetive,
		},
		FieldTypeSliceOfStruct: {
			FieldTypeSliceOfStruct: genStructSliceStructSliceMapping,
		},
		FieldTypeSliceOfPrimetive: {
			FieldTypeSliceOfPrimetive: genPrimetivePrimetive,
		},
		FieldTypePointerToPrimetive: {
			FieldTypePrimetive:          genPtrPrimetivePrimetive,
			FieldTypePointerToPrimetive: genPrimetivePrimetive,
		},
		FieldTypePointerToStruct: {
			FieldTypePointerToStruct: genPtrStructPtrStructMapping,
			FieldTypeStruct:          genPtrStructStructMapping,
		},
	}

	f.Comment("Code generated by go-mapping-jam. DO NOT EDIT.")

	for _, mapFunc := range mapFuncs {
		gm := generatedMapper{
			generatedFile:    &g,
			name:             mapFunc.Name,
			from:             mapFunc.Source,
			to:               mapFunc.Target,
			rules:            mapFunc.Rules,
			fieldGenMapFuncs: fieldGenMap,
		}
		gm.generateMapFunc()
	}

	err := f.Save(outputFile)
	if err != nil {
		log.Fatalf("failed to save file: %v", err)
	}
}

func (gm *generatedMapper) initSignature() {
	gm.Statement = gm.Func().Id(gm.name)
	if gm.from.Path == currentPath {
		gm.Statement.Params(jen.Id("src").Id(gm.from.Name))
	} else {
		gm.Statement.Params(jen.Id("src").Qual(gm.from.Path, gm.from.Name))
	}

	if gm.to.Path == currentPath {
		gm.Statement.Id(gm.to.Name)
	} else {
		gm.Statement.Qual(gm.to.Path, gm.to.Name)
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
