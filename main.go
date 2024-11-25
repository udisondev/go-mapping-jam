package main

const projectName = "github.com/udisondev/go-mapping-jam"

var currentPath =  projectName + "/mapping"

func main() {
	registerRuleParser("qual", parseQualRule)
	registerRuleParser("enum", parseEnumRule)

	mapFuncs := parse("./mapper/mapper.go")
	
	outputFile := "./mapper/mapper_impl.go"
	generateCodeWithJennifer(outputFile, mapFuncs)

}
