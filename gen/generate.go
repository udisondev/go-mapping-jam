package gen

import (
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/dave/jennifer/jen"
	"github.com/udisondev/go-mapping-jam/mapp"
)

const charset = "abcdefghijklmnopqrstuvwxyz"

type mapperFunc struct {
	generatedFn        *jen.Statement
	mapper             mapp.Mapper
	submappers         map[string]string
	fieldMapGenerators map[mapp.TypeFamily]map[mapp.TypeFamily]func(bl mapperBlock, s, t mapp.Field)
}

type submapperFunc struct {
	generatedFn        *jen.Statement
	mapper             mapp.Mapper
	source             *mapp.Field
	target             *mapp.Field
	submappers         map[string]string
	fieldMapGenerators map[mapp.TypeFamily]map[mapp.TypeFamily]func(bl mapperBlock, s, t mapp.Field)
}

type mapperBlock struct {
	*jen.Group
	mapperFunc
}

func Generate(mf mapp.MapperFile) {
	f := jen.NewFile("mapper")

	f.Comment("Code generated by go-mapping-jam. DO NOT EDIT.")

	submappers := make(map[string]string)
	fieldMapGenerators := map[mapp.TypeFamily]map[mapp.TypeFamily]func(bl mapperBlock, s, t mapp.Field){
		mapp.FieldTypeBasic: {
			mapp.FieldTypeBasic:   generateBasicToBasic,
			mapp.FieldTypePointer: generateBasicToPointer},
		mapp.FieldTypePointer: {mapp.FieldTypeBasic: generatePointerToBasic},
	}
	for _, m := range mf.Mappers() {
		mfn := mapperFunc{
			generatedFn:        f.Func().Id(m.Name()),
			mapper:             m,
			submappers:         submappers,
			fieldMapGenerators: fieldMapGenerators,
		}
		mfn.generateSignature()
		mfn.generateBlock()
	}

	err := f.Save("./mapper/mapper_impl.go")
	if err != nil {
		log.Fatalf("failed to save file: %v", err)
	}
}

func (m mapperFunc) generateBlock() {
	m.generatedFn.BlockFunc(func(g *jen.Group) {
		bl := mapperBlock{
			Group:      g,
			mapperFunc: m,
		}
		target := m.mapper.Target()
		_, t := target.Type()
		bl.Id("target").Op(":=").Qual(target.Path(), t+"{}")
		for _, tf := range target.Fields() {
			bl.generateTargetMapping(tf)
		}
		bl.Return(jen.Id("target"))
	})

}

func (bl mapperBlock) generateTargetMapping(target mapp.Field) {
	_, ok := bl.mapper.RulesByFieldFullNameAndType(target.FullName(), mapp.RuleTypeIgnore)
	if ok {
		return
	}

	source, ok := bl.mapper.SourceFieldByTarget(target.FullName())
	if !ok {
		panic(fmt.Sprintf("not found source field for target: %s", target.FullName()))
	}

	genFn, ok := bl.fieldMapGenerators[source.Type().TypeFamily()][target.Type().TypeFamily()]
	if !ok {
		return
	}

	genFn(bl, source, target)
}

func (m mapperFunc) generateSignature() {
	for i, p := range m.mapper.Params() {
		_, typeName := p.Type()
		pname := p.Name()
		if i == 0 {
			pname = "src"
		}
		m.generatedFn.Params(jen.Id(pname).Qual(p.Path(), typeName))
	}

	for _, r := range m.mapper.Results() {
		_, typeName := r.Type()
		m.generatedFn.Qual(r.Path(), typeName)
	}
}

func generateBasicToBasic(bl mapperBlock, s, t mapp.Field) {
	bl.Id("target").Dot(t.Name()).Op("=").Id("src").Dot(s.Name())
}

func generatePointerToBasic(bl mapperBlock, s, t mapp.Field) {
	pt, ok := s.Type().(mapp.PointerType)
	if !ok {
		panic("is not a pointer")
	}

	if t.Type().Type() != s.Type().Type() {
		panic(fmt.Sprintf(
			"could not mapp different types source: '*%s' target: %s",
			pt.Elem().TypeFamily(),
			t.Type().TypeFamily()))
	}

	bl.If(
		jen.Id("src").Dot(s.Name()).Op("!=").Nil(),
	).Block(
		jen.Id("target").Dot(t.Name()).Op("=").Add(jen.Op("*")).Id("src").Dot(s.Name()),
	)
}

func generateBasicToPointer(bl mapperBlock, s, t mapp.Field) {
	pt, ok := t.Type().(mapp.PointerType)
	if !ok {
		panic("is not a pointer")
	}

	if pt.Elem().TypeFamily() != mapp.FieldTypeBasic {
		panic("source refers to not basic")
	}

	if s.Type().TypeFamily() != pt.Elem().TypeFamily() {
		panic(fmt.Sprintf(
			"could not mapp different types source: '%s' target: pointer to %s",
			s.Type().TypeFamily(),
			pt.Elem().TypeFamily()))
	}

	bl.Id("target").Dot(t.Name()).Op("=").Add(jen.Op("&")).Id("src").Dot(s.Name())
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
