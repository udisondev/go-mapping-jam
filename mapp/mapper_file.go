package mapp

import (
	"go/ast"
	"strings"
)

type MapperFile struct {
	spec *ast.File
}

func NewMapperFile(node *ast.File) *MapperFile {
	return &MapperFile{
		spec: node,
	}
}

func (mf MapperFile) Imports() []Import {
	imports := make([]Import, 0, len(mf.spec.Imports))
	for _, i := range mf.spec.Imports {
		imp := Import{
			spec: i,
		}
		imports = append(imports, imp)
	}

	return imports
}

func (mf MapperFile) Mappers() []Mapper {
	methodList := make([]Mapper, 0)
	ast.Inspect(mf.spec, func(n ast.Node) bool {
		iface, ok := n.(*ast.InterfaceType)
		if !ok {
			return true
		}

		imports := mf.Imports()
	searchMethodLoop:
		for _, v := range iface.Methods.List {
			for _, d := range v.Doc.List {
				if strings.Contains(d.Text, "@emapper") {
					continue searchMethodLoop
				}
			}

			methodList = append(methodList, Mapper{
				spec:    v,
				imports: imports,
			})
		}

		return false
	})

	return methodList
}
