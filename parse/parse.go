package parse

import (
	"fmt"
	"go/parser"
	"go/token"

	"github.com/udisondev/go-mapping-jam/mapp"
)

func File(filePath string) *mapp.MapperFile {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		panic(fmt.Sprintf("failed to parse file: %v", err))
	}

	return mapp.NewMapperFile(node)
}
