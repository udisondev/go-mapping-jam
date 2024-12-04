package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/udisondev/go-mapping-jam/mapp"
)

func structToStruct(bl mapperBlock, s, t mapp.Field) {
	hash := fieldsHash(s, t)
	submapperName, submapperExists := bl.submappers[hash]
	if !submapperExists {
		submapperName = genRandomName(10)
		bl.submappers[hash] = submapperName
	}
	if s.Type().TypeFamily() != mapp.FieldTypePointer {
		bl.Id("target").Dot(t.Name()).Op("=").Id(submapperName).Call(jen.Id("src").Dot(s.Name()))
	} else {
		bl.Id("target").Dot(t.Name()).Op("=").Id(submapperName).Call(jen.Add(jen.Op("*")).Id("src").Dot(s.Name()))
	}

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