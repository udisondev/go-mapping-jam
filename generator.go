package main

import (
	"log"
	"math/rand/v2"
	"time"

	jen "github.com/dave/jennifer/jen"
)


func generateCodeWithJennifer(outputFile string, mapFuncs map[string]MapFunc) {
	f := jen.NewFile("main")
	f.Comment("Code generated by go-mapping-jam. DO NOT EDIT.")

	// Создаём структуру
	f.Type().Id("MapperImpl").Struct()

	// Генерируем функции
	for _, mapFunc := range mapFuncs {
		generateMapFunc(f, mapFunc)
	}

	// Сохраняем файл
	err := f.Save(outputFile)
	if err != nil {
		log.Fatalf("failed to save file: %v", err)
	}
}

func generateMapFunc(f *jen.File, mapFunc MapFunc) {
	sourceStruct := mapFunc.Mappable[Source]
	targetStruct := mapFunc.Mappable[Target]

	// Основной метод
	mapF := f.Func().
		Params(jen.Id("m").Op("*").Id("MapperImpl")).
		Id(mapFunc.Name)
	if sourceStruct.Pack == thisPack {
		mapF.Params(jen.Id("src").Id(sourceStruct.FullType.StructName))
	} else {
		mapF.Params(jen.Id("src").Qual(sourceStruct.Pack.Path, sourceStruct.FullType.StructName))
	}

	if targetStruct.Pack == thisPack {
		mapF.Id(targetStruct.FullType.StructName)
	} else {
		mapF.Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName)
	}

	mapF.BlockFunc(func(g *jen.Group) {
		if targetStruct.Pack == thisPack {
			g.Id("target").Op(":=").Id(targetStruct.FullType.StructName + "{}")
		} else {
			g.Id("target").Op(":=").Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName+"{}")
		}

		for targetFieldName, targetField := range mapFunc.Mappable[Target].Fields {
			qual := mapFunc.Rules[Qual]
			sourceFieldName := targetFieldName
			if q, ok := qual.(QualRule); ok && q.TargetName == targetFieldName {
				sourceFieldName = q.SourceName
			}

			if sourceField, ok := mapFunc.Mappable[Source].Fields[sourceFieldName]; ok {
				// Если это примитивное поле
				if _, ok := targetField.(*Primetive); ok {
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("src").Dot(string(sourceFieldName))
				} else if targetStruct, ok := targetField.(*Struct); ok {
					hash := string(targetStruct.Hash()) + string(sourceField.(*Struct).Hash())
					subMethodName, ok := subMappers[hash]
					if !ok {
						subMethodName = genRandomName(15)
						subMappers[hash] = subMethodName
					}
					// Если это структура
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("m").Dot(subMethodName).Call(jen.Id("src").Dot(string(sourceFieldName)))

					if !ok {
						// Генерируем подметод для вложенных структур
						generateSubMapper(f, subMethodName, sourceField.(*Struct), targetStruct, mapFunc)
					}
				}
			}
		}

		g.Return(jen.Id("target"))
	})
}

var subMappers = make(map[string]string)

func generateSubMapper(f *jen.File, methodName string, sourceStruct *Struct, targetStruct *Struct, mapFunc MapFunc) {
	// Основной метод
	mapF := f.Func().
		Params(jen.Id("m").Op("*").Id("MapperImpl")).
		Id(methodName)
	if sourceStruct.Pack == thisPack {
		mapF.Params(jen.Id("src").Id(sourceStruct.FullType.StructName))
	} else {
		mapF.Params(jen.Id("src").Qual(sourceStruct.Pack.Path, sourceStruct.FullType.StructName))
	}

	if targetStruct.Pack == thisPack {
		mapF.Id(targetStruct.FullType.StructName)
	} else {
		mapF.Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName)
	}

	mapF.BlockFunc(func(g *jen.Group) {
		if targetStruct.Pack == thisPack {
			g.Id("target").Op(":=").Id(targetStruct.FullType.StructName + "{}")
		} else {
			g.Id("target").Op(":=").Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName+"{}")
		}

		for targetFieldName, targetField := range targetStruct.Fields {
			qual := mapFunc.Rules[Qual]
			sourceFieldName := targetFieldName
			if q, ok := qual.(QualRule); ok && q.TargetName == targetFieldName {
				sourceFieldName = q.SourceName
			}

			if sourceField, ok := sourceStruct.Fields[sourceFieldName]; ok {
				if _, ok := sourceField.(*Primetive); ok {
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("src").Dot(string(sourceFieldName))
				} else if nestedSourceStruct, ok := sourceField.(*Struct); ok {
					hash := string(nestedSourceStruct.Hash()) + string(targetField.(*Struct).Hash())
					methodName, ok := subMappers[hash]
					if !ok {
						methodName = genRandomName(15)
						subMappers[hash] = methodName
					}
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("m").Dot(methodName).Call(jen.Id("src").Dot(string(sourceFieldName)))

					if !ok {
						// Рекурсивно генерируем вложенные подметоды
						generateSubMapper(f, methodName, nestedSourceStruct, targetField.(*Struct), mapFunc)
					}
				}
			}
		}

		g.Return(jen.Id("target"))
	})
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func genRandomName(length int) string {
	seed := time.Now().UnixNano()

	src := rand.NewPCG(uint64(seed), uint64(seed>>32))
	r := rand.New(src)

	// Буфер для сборки строки
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.IntN(len(charset))]
	}
	return string(result)
}