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

	sourceField, ok := mf.source().Fields[sourceFieldName]
	if !ok {
		log.Fatalf("source field not found for target: %s", mf.field.FullName())
	}

	switch tarT := mf.field.Desc.(type) {
	case *Primetive:
		switch srcT := sourceField.Desc.(type) {
		case *Primetive:
			mf.genPrimetivePrimetive(sourceFieldName)
		case *Pointer:
			_, ok := srcT.To.(*Primetive)
			if !ok {
				panic(fmt.Sprintf("source field is not a pointer to primetive for target: %s", mf.field.FullName()))
			}
			mf.genPrimetivePtrPrimetive(sourceFieldName)
		}
	case *Struct:
		switch srcT := sourceField.Desc.(type) {
		case *Struct:
			mf.genStructStructMapping(sourceFieldName, sourceField)
		case *Pointer:
			_, ok := srcT.To.(*Struct)
			if !ok {
				panic(fmt.Sprintf("source field is not a pointer to struct for target: %s", mf.field.FullName()))
			}
			mf.genStructPtrStructMapping(sourceFieldName, sourceField)
		}
	case *Enum:
		switch srcT := sourceField.Desc.(type) {
		case *Enum:
			if srcT.Name == tarT.Name {
				mf.genPrimetivePrimetive(sourceFieldName)
				break
			}
			mf.genEnumMapping(sourceFieldName, sourceField)
		}
	case *Slice:
		switch tarT.Of.(type) {
		case *Primetive:
			slc, ok := sourceField.Desc.(*Slice)
			if !ok {
				panic(fmt.Sprintf("source field is not a slice for target: %s", mf.field.FullName()))
			}
			_, ok = slc.Of.(*Primetive)
			if !ok {
				panic(fmt.Sprintf("source field is not a slice of primetive for target: %s", mf.field.FullName()))
			}
			mf.group().Id("target").Dot(mf.name).Op("=").Id("src").Dot(sourceFieldName)
		case *Struct:
			slc, ok := sourceField.Desc.(*Slice)
			if !ok {
				panic(fmt.Sprintf("source field is not a slice for target: %s", mf.field.FullName()))
			}
			_, ok = slc.Of.(*Struct)
			if !ok {
				panic(fmt.Sprintf("source field is not a slice of struct for target: %s", mf.field.FullName()))
			}
			mf.genStructSliceStructSliceMapping(sourceFieldName, sourceField)
		}
	case *Pointer:
		switch tarT.To.(type) {
		case *Primetive:
			switch srcT := sourceField.Desc.(type) {
			case *Primetive:
				mf.genPtrPrimetivePrimetive(sourceFieldName)
			case *Pointer:
				_, ok := srcT.To.(*Primetive)
				if !ok {
					panic(fmt.Sprintf("source field is not a pointer to primetive for target: %s", mf.field.FullName()))
				}
				mf.genPrimetivePrimetive(sourceFieldName)
			}
		case *Struct:
			switch srcT := sourceField.Desc.(type) {
			case *Struct:
				mf.genPtrStructStructMapping(sourceFieldName, sourceField)
			case *Pointer:
				_, ok := srcT.To.(*Struct)
				if !ok {
					panic(fmt.Sprintf("source field is not a pointer to struct for target: %s", mf.field.FullName()))
				}
				mf.genPtrStructPtrStructMapping(sourceFieldName, sourceField)
			}
		}
	}

}

func (mf *mappedField) genEnumMapping(sourceFieldName string, sourceField *Field) {

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

func (mf *mappedField) genPrimetivePrimetive(sourceFieldName string) {
	mf.group().Id("target").Dot(mf.name).Op("=").Id("src").Dot(sourceFieldName)
}

func (mf *mappedField) genPtrPrimetivePrimetive(sourceFieldName string) {
	mf.group().Id("target").Dot(mf.name).Op("=").Add(jen.Op("&")).Id("src").Dot(sourceFieldName)
}

func (mf *mappedField) genPtrStructPtrStructMapping(sourceFieldName string, sourceField *Field) {
	nestedSourceStruct, ok := sourceField.Desc.(*Pointer).To.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	targetField, ok := mf.field.Desc.(*Pointer).To.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	hash := nestedSourceStruct.Hash() + targetField.Hash()
	methodName, ok := mf.submappers()[hash]
	if !ok {
		methodName = genRandomName(15)
		mf.submappers()[hash] = methodName
	}

	mf.group().
		If(
			jen.Id("src").Dot(sourceFieldName).Op("!=").Nil(),
		).
		Block(
			jen.Id(methodName+"Result").Op(":=").Id(methodName).Call(jen.Add(jen.Op("*").Id("src").Dot(sourceFieldName))),
			jen.Id("target").Dot(mf.field.Name).Op("=").Add(jen.Op("&")).Id(methodName+"Result"),
		)

	if !ok {
		sbm := generatedMapper{
			name:       methodName,
			from:       nestedSourceStruct,
			to:         targetField,
			isFromPrt:  true,
			isToPtr:    true,
			file:       mf.file(),
			rules:      mf.rules(),
			submappers: mf.submappers(),
		}
		sbm.generateMapFunc()
	}
}

func (mf *mappedField) genPtrStructStructMapping(sourceFieldName string, sourceField *Field) {
	nestedSourceStruct, ok := sourceField.Desc.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	targetField, ok := mf.field.Desc.(*Pointer).To.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	hash := nestedSourceStruct.Hash() + targetField.Hash()
	methodName, ok := mf.submappers()[hash]
	if !ok {
		methodName = genRandomName(15)
		mf.submappers()[hash] = methodName
	}

	mf.group().Id(methodName + "Result").Op(":=").Id(methodName).Call(jen.Id("src").Dot(sourceFieldName))
	mf.group().Id("target").Dot(mf.field.Name).Op("=").Add(jen.Op("&")).Id(methodName + "Result")

	if !ok {
		sbm := generatedMapper{
			name:       methodName,
			from:       nestedSourceStruct,
			to:         targetField,
			isFromPrt:  true,
			isToPtr:    true,
			file:       mf.file(),
			rules:      mf.rules(),
			submappers: mf.submappers(),
		}
		sbm.generateMapFunc()
	}
}

func (mf *mappedField) genStructPtrStructMapping(sourceFieldName string, sourceField *Field) {
	nestedSourceStruct, ok := sourceField.Desc.(*Pointer).To.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	targetField, ok := mf.field.Desc.(*Struct)
	if !ok {
		panic("is not a struct")
	}

	hash := nestedSourceStruct.Hash() + targetField.Hash()
	methodName, ok := mf.submappers()[hash]
	if !ok {
		methodName = genRandomName(15)
		mf.submappers()[hash] = methodName
	}

	mf.group().
		If(
			jen.Id("src").Dot(sourceFieldName).Op("!=").Nil(),
		).
		Block(
			jen.Id("target").Dot(mf.field.Name).Op("=").Id(methodName).Call(jen.Add(jen.Op("*").Id("src").Dot(sourceFieldName))),
		)

	if !ok {
		sbm := generatedMapper{
			name:       methodName,
			from:       nestedSourceStruct,
			to:         targetField,
			isFromPrt:  true,
			isToPtr:    true,
			file:       mf.file(),
			rules:      mf.rules(),
			submappers: mf.submappers(),
		}
		sbm.generateMapFunc()
	}
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

func (mf *mappedField) genStructSliceStructSliceMapping(sourceFieldName string, sourceField *Field) {
	targetStruct := mf.field.Desc.(*Slice).Of.(*Struct)

	nestedSourceStruct := sourceField.Desc.(*Slice).Of.(*Struct)
	hash := nestedSourceStruct.Hash() + targetStruct.Hash()
	methodName, ok := mf.submappers()[hash]
	if !ok {
		methodName = genRandomName(15)
		mf.submappers()[hash] = methodName
	}
	mf.group().Id("target"+mf.name+"Slice").Op(":=").Make(jen.Index().Qual(targetStruct.Path, targetStruct.Name), jen.Lit(0), jen.Len(jen.Id("target").Dot(mf.name)))
	mf.group().
		For(
			jen.List(jen.Id("_"), jen.Id("it")).Op(":=").Range().Id("src").Dot(sourceFieldName),
		).
		Block(
			jen.Id("target"+mf.name+"Slice").Op("=").Append(jen.Id("target"+mf.name+"Slice"), jen.Id(methodName).Call(jen.Id("it"))),
		)

	mf.group().Id("target").Dot(mf.name).Op("=").Id("target" + mf.name + "Slice")

	if !ok {
		sbm := generatedMapper{
			name:       methodName,
			from:       nestedSourceStruct,
			to:         targetStruct,
			file:       mf.file(),
			rules:      mf.rules(),
			submappers: mf.submappers(),
		}
		sbm.generateMapFunc()
	}
}

func (mf *mappedField) findQualRule() (QualRule, bool) {
	for _, v := range mf.rules()[RuleTypeQual] {
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
				source:     func() *Struct { return gm.from },
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
		return buildMappingCase(tarPtr.To.fieldType(), sourceField.Desc.fieldType(), true, false)
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
		return buildMappingCase(tarPtr.To.fieldType(), srcPtr.To.fieldType(), true, true)
	}

	if gf.field.Desc.fieldType() != PointerType && sourceField.Desc.fieldType() == PointerType {
		srcPtr, ok := sourceField.Desc.(*Pointer)
		if !ok {
			panic("is not a pointer")
		}
		return buildMappingCase(gf.field.Desc.fieldType(), srcPtr.To.fieldType(), false, true)
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
