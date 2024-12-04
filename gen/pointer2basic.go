package gen

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/udisondev/go-mapping-jam/mapp"
)

func pointerToBasic(bl mapperBlock, s, t mapp.Field) {
	pt, ok := s.Type().(mapp.PointerType)
	if !ok {
		panic("is not a pointer")
	}

	if t.Type().TypeName() != s.Type().TypeName() {
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