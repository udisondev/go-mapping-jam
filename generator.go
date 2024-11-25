package main

import (
	"log"
	"math/rand/v2"
	"time"

	jen "github.com/dave/jennifer/jen"
)


func generateCodeWithJennifer(outputFile string, mapFuncs map[string]MapFunc) {
	f := jen.NewFile("mapper")
	f.Comment("Code generated by go-mapping-jam. DO NOT EDIT.")

	f.Type().Id("MapperImpl").Struct()

	for _, mapFunc := range mapFuncs {
		generateMapFunc(f, mapFunc)
	}

	err := f.Save(outputFile)
	if err != nil {
		log.Fatalf("failed to save file: %v", err)
	}
}

func generateMapFunc(f *jen.File, mapFunc MapFunc) {
	sourceStruct := mapFunc.Source
	targetStruct := mapFunc.Target

	mapF := f.Func().
		Params(jen.Id("m").Op("*").Id("MapperImpl")).
		Id(mapFunc.Name)
	if sourceStruct.Path == currentPath {
		mapF.Params(jen.Id("src").Id(sourceStruct.Name))
	} else {
		mapF.Params(jen.Id("src").Qual(sourceStruct.Path, sourceStruct.Name))
	}

	if targetStruct.Path == currentPath {
		mapF.Id(targetStruct.Name)
	} else {
		mapF.Qual(targetStruct.Path, targetStruct.Name)
	}

	mapF.BlockFunc(func(g *jen.Group) {
		if targetStruct.Path == currentPath {
			g.Id("target").Op(":=").Id(targetStruct.Name + "{}")
		} else {
			g.Id("target").Op(":=").Qual(targetStruct.Path, targetStruct.Name+"{}")
		}

		for targetFieldName, targetField := range mapFunc.Target.Fields {
			quals := mapFunc.Rules[Qual]
			sourceFieldName := targetFieldName
			var mname string
			for _, q := range quals {
				if q, ok := q.(QualRule); ok && q.TargetName == targetField.FullName() {
					sourceFieldName = q.SourceName
					if q.MethodName != "" {
						mname = q.MethodName
					}
				}	
			}
			
			if sourceField, ok := mapFunc.Source.Fields[sourceFieldName]; ok {
				if mname != "" {
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id(mname).Call(jen.Id("src").Dot(sourceFieldName))
				} else if _, ok := targetField.Desc.(*Primetive); ok {
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("src").Dot(string(sourceFieldName))
				} else if targetStruct, ok := targetField.Desc.(*Struct); ok {
					hash := string(targetStruct.Hash()) + string(sourceField.Desc.(*Struct).Hash())
					subMethodName, ok := subMappers[hash]
					if !ok {
						subMethodName = genRandomName(15)
						subMappers[hash] = subMethodName
					}
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("m").Dot(subMethodName).Call(jen.Id("src").Dot(sourceFieldName))

					if !ok {
						generateSubMapper(f, subMethodName, sourceField.Desc.(*Struct), targetStruct, mapFunc)
					}
				}
			}
		}

		g.Return(jen.Id("target"))
	})
}

var subMappers = make(map[string]string)

func generateSubMapper(f *jen.File, methodName string, sourceStruct *Struct, targetStruct *Struct, mapFunc MapFunc) {
	mapF := f.Func().
		Params(jen.Id("m").Op("*").Id("MapperImpl")).
		Id(methodName)
	if sourceStruct.Path == currentPath {
		mapF.Params(jen.Id("src").Id(sourceStruct.Name))
	} else {
		mapF.Params(jen.Id("src").Qual(sourceStruct.Path, sourceStruct.Name))
	}

	if targetStruct.Path == currentPath {
		mapF.Id(targetStruct.Name)
	} else {
		mapF.Qual(targetStruct.Path, targetStruct.Name)
	}

	mapF.BlockFunc(func(g *jen.Group) {
		if targetStruct.Path == currentPath {
			g.Id("target").Op(":=").Id(targetStruct.Name + "{}")
		} else {
			g.Id("target").Op(":=").Qual(targetStruct.Path, targetStruct.Name+"{}")
		}

		for targetFieldName, targetField := range targetStruct.Fields {
			quals := mapFunc.Rules[Qual]
			sourceFieldName := targetFieldName
			var mname string
			for _, q := range quals {
				log.Printf("target full name: %s", targetField.FullName())
				if q, ok := q.(QualRule); ok && q.TargetName == targetField.FullName() {
					sourceFieldName = q.SourceName
					if q.MethodName != "" {
						mname = q.MethodName
					}
				}	
			}

			if sourceField, ok := sourceStruct.Fields[sourceFieldName]; ok {
				if mname != "" {
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id(mname).Call(jen.Id("src").Dot(sourceFieldName))
				} else if _, ok := sourceField.Desc.(*Primetive); ok {
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("src").Dot(string(sourceFieldName))
				} else if nestedSourceStruct, ok := sourceField.Desc.(*Struct); ok {
					hash := string(nestedSourceStruct.Hash()) + string(targetField.Desc.(*Struct).Hash())
					methodName, ok := subMappers[hash]
					if !ok {
						methodName = genRandomName(15)
						subMappers[hash] = methodName
					}
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("m").Dot(methodName).Call(jen.Id("src").Dot(string(sourceFieldName)))

					if !ok {
						generateSubMapper(f, methodName, nestedSourceStruct, targetField.Desc.(*Struct), mapFunc)
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

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.IntN(len(charset))]
	}
	return string(result)
}