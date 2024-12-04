package main

import (
	"fmt"

	"github.com/udisondev/go-mapping-jam/gen"
	"github.com/udisondev/go-mapping-jam/mapp"
	"github.com/udisondev/go-mapping-jam/parse"
)

const ProjectName = "github.com/udisondev/go-mapping-jam"

var CurrentPath = ProjectName + "/mapping"

func main() {
	// rule.RegisterRuleParser("qual", parseQualRule)
	// rule.RegisterRuleParser("enum", parseEnumRule)

	mapperFile := parse.File("./mapper/mapper.go")
	gen.Generate(mapperFile)
	// check(mapperFile)
}

func check(mapperFile mapp.MapperFile) {
	imports := mapperFile.Imports()
	for _, i := range imports {
		fmt.Println("Import", "Alias:", i.Alias(), "Path:", i.Path())
	}
	println("----------------------------")
	mappers := mapperFile.Mappers()
	for _, m := range mappers {
		f, ok := m.SourceFieldByTarget("t.Profile.Phone")
		if ok {
			fmt.Printf("found f: %s\n", f.FullName())
		}
		fmt.Printf("mapper.Name(): %v\n", m.Name())
		params := m.Params()
		for _, p := range params {
			fmt.Printf("param.Name(): %v\n", p.Name())
			pack, t := p.Type()
			fmt.Println("param.Type():", pack, t)
			fmt.Printf("param.Path(): %v\n", p.Path())
		}
		comments := m.Comments()
		for _, c := range comments {
			fmt.Printf("comment.Value(): %v\n", c.Value())
		}

		for _, r := range m.Rules() {
			fmt.Printf("rule: %v\n", r.Value())
		}

		for _, r := range m.Results() {
			fmt.Printf("rule.Name(): %v\n", r.Name())
			pack, t := r.Type()
			fmt.Println("rule.Type():", pack, t)
			fmt.Printf("rule.Path(): %v\n", r.Path())
		}

		src := m.Source()
		sfields := src.Fields()
		fmt.Println("Source fields: ")
		for _, f := range sfields {
			deepFields(f)
		}
		println("----------------------------")

		target := m.Target()
		tfields := target.Fields()
		fmt.Println("target fields: ")
		for _, f := range tfields {
			deepFields(f)
		}
		println("----------------------------")
	}
}

func deepFields(f mapp.Field) {
	fmt.Printf("inner field: %s typeFamily: %s type: %s path: %s\n", f.FullName(), f.Type().TypeFamily(), f.Type().TypeName(), f.Type().Path())
	fields := f.Fields()
	if len(fields) == 0 {
		return
	}
	for _, ff := range fields {
		deepFields(ff)
	}
}
