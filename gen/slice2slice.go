package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/udisondev/go-mapping-jam/mapp"
)

func sliceToSlice(bl mapperBlock, s, t mapp.Field) {

	sslice, ok := s.Type().(mapp.SliceType)
	if !ok {
		panic("is not a slice")
	}

	tslice, ok := t.Type().(mapp.SliceType)
	if !ok {
		panic("is not a slice")
	}

	switch {
		case sslice.Elem().TypeFamily() == mapp.FieldTypeBasic &&
		tslice.Elem().TypeFamily() == mapp.FieldTypeBasic:
		basicToBasic(bl, s, t)
	case sslice.Elem().TypeFamily() == mapp.FieldTypeStruct &&
		tslice.Elem().TypeFamily() == mapp.FieldTypeStruct:
		hash := fieldsHash(s, t)
		submapperName, submapperExists := bl.submappers[hash]
		if !submapperExists {
			submapperName = genRandomName(10)
			bl.submappers[hash] = submapperName
		}
	
		targetSliceName := "target" + t.Name() + "Slice"
		targetTypePath := t.Type().Path()
		bl.Id(targetSliceName).
			Op(":=").
			Make(
				jen.Index().Qual(targetTypePath, t.Type().TypeName()),
				jen.Lit(0),
				jen.Len(jen.Id("target").Dot(t.Name())))
		bl.
			For(
				jen.List(jen.Id("_"), jen.Id("it")).Op(":=").Range().Id("src").Dot(s.Name()),
			).
			Block(
				jen.Id(targetSliceName).Op("=").Append(jen.Id(targetSliceName), jen.Id(submapperName).Call(jen.Id("it"))),
			)
	
		bl.Id("target").Dot(t.Name()).Op("=").Id(targetSliceName)
	
		if !submapperExists {
			mfn := mapperFunc{
				generatedFn:        bl.file.Func().Id(submapperName),
				mapper:             bl.mapper,
				file:               bl.mapperFunc.file,
				source:             s,
				target:             t,
				submappers:         bl.submappers,
				fieldMapGenerators: bl.fieldMapGenerators,
			}
			bl.file.Line()
			mfn.generateSignature()
			mfn.generateBlock()
		}
	}

	
}