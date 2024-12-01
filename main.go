package main

import (
	"fmt"

	"github.com/udisondev/go-mapping-jam/mapp"
	"github.com/udisondev/go-mapping-jam/parse"
)

const ProjectName = "github.com/udisondev/go-mapping-jam"

var CurrentPath = ProjectName + "/mapping"

func main() {
	// rule.RegisterRuleParser("qual", parseQualRule)
	// rule.RegisterRuleParser("enum", parseEnumRule)

	mapperFile := parse.File("./mapper/mapper.go")
	check(mapperFile)
}

func check(mapperFile *mapp.MapperFile) {
	imports := mapperFile.Imports()
	for _, i := range imports {
		fmt.Println("Import", "Alias:", i.Alias(), "Path:", i.Path())
	}
	println("----------------------------")
	mappers := mapperFile.Mappers()
	for _, m := range mappers {
		fmt.Printf("m.Name(): %v\n", m.Name())
		params := m.Params()
		for _, p := range params {
			fmt.Printf("p.Name(): %v\n", p.Name())
			pack, t := p.Type()
			fmt.Println("p.Type():", pack, t)
			fmt.Printf("p.Path(): %v\n", p.Path())
		}
		comments := m.Comments()
		for _, c := range comments {
			fmt.Printf("c.Value(): %v\n", c.Value())
		}

		for _, r := range m.Rules() {
			fmt.Printf("rule: %v\n", r.Value())
		}

		for _, r := range m.Results() {
			fmt.Printf("r.Name(): %v\n", r.Name())
			pack, t := r.Type()
			fmt.Println("r.Type():", pack, t)
			fmt.Printf("r.Path(): %v\n", r.Path())
		}
		println("----------------------------")
	}
}
