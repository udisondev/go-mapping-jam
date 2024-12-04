package gen

import (
	"github.com/dave/jennifer/jen"
	"github.com/udisondev/go-mapping-jam/mapp"
)

func pointerToStruct(bl mapperBlock, s, t mapp.Field) {
	bl.If(
		jen.Id("src").Dot(s.Name()).Op("!=").Nil(),
	).BlockFunc(
		func(g *jen.Group) {
			bl.Group = g
			structToStruct(bl, s, t)
		},
	)
}