package gen

import (
	"github.com/udisondev/go-mapping-jam/mapp"
)

func pointerToPointer(bl mapperBlock, s, t mapp.Field) {
	// TODO доделать
	// spt, ok := s.Type().(mapp.PointerType)
	// if !ok {
	// 	panic("is not a pointer")
	// }
	// tpt, ok := t.Type().(mapp.PointerType)
	// if !ok {
	// 	panic("is not a pointer")
	// }

	// switch {
	// case spt.Elem().TypeFamily() == mapp.FieldTypeBasic && tpt.Elem().TypeFamily() == mapp.FieldTypeBasic:
	// 	basicToBasic(bl, s, t)
	// }
	// bl.If(
	// 	jen.Id("src").Dot(s.Name()).Op("!=").Nil(),
	// ).BlockFunc(
	// 	func(g *jen.Group) {
	// 		bl.Group = g
	// 		structToStruct(bl, s, t)
	// 	},
	// )
}