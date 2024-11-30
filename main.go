package main

import "github.com/udisondev/go-mapping-jam/rule"

const projectName = "github.com/udisondev/go-mapping-jam"

var currentPath =  projectName + "/mapping"

func main() {
	rule.RegisterRuleParser("qual", parseQualRule)
	rule.RegisterRuleParser("enum", parseEnumRule)

	mapFuncs := parse("./mapper/mapper.go")
	
	outputFile := "./mapper/mapper_impl.go"
	generateCodeWithJennifer(outputFile, mapFuncs)
}
