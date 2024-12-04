package gen

import "github.com/udisondev/go-mapping-jam/mapp"

func basicToBasic(bl mapperBlock, s, t mapp.Field) {
	bl.Id("target").Dot(t.Name()).Op("=").Id("src").Dot(s.Name())
}