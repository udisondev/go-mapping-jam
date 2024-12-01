package main

// import (
// 	"go/ast"
// 	"go/parser"
// 	"go/token"
// 	"log"
// 	"path/filepath"
// 	"regexp"
// 	"strings"

// 	. "github.com/udisondev/go-mapping-jam/mapp"
// 	"github.com/udisondev/go-mapping-jam/rule"
// 	"golang.org/x/tools/go/packages"
// )

// var pkgs = make(map[string]*packages.Package)

// func parse(filePath string) map[string]Mapper {
	

// 	mapperImports := make([]struct{ alias, path string }, 0, len(node.Imports))
// 	for _, v := range node.Imports {
// 		mapperImports = append(mapperImports, extractMapImport(v))
// 	}

// 	mappersMap := make(map[string]Mapper)

// 	for _, v := range extractMethods(node) {
// 		mappingRules := make(map[rule.Type][]rule.Any)
// 		if v.Doc != nil {
// 			for _, mpr := range v.Doc.List {
// 				rules := parseRules(strings.TrimSpace(strings.ReplaceAll(mpr.Text, "//", "")))
// 				for _, r := range rules {
// 					mappingRules[r.Type()] = append(mappingRules[r.Type()], r)
// 				}
// 			}
// 		}

// 		if fType, ok := v.Type.(*ast.FuncType); ok {
// 			source := extractMappingRoot(fType.Params, mapperImports)
// 			target := extractMappingRoot(fType.Results, mapperImports)
// 			m := Mapper{
// 				Name:   v.Names[0].Name,
// 				Source: source,
// 				Target: target,
// 				Rules:  mappingRules,
// 			}

// 			mappersMap[m.Name] = m
// 		}

// 	}

// 	cfg := &packages.Config{
// 		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
// 		Fset: fset,
// 	}

// 	for _, v := range mappersMap {
// 		packFunc := func(dir string) *packages.Package {
// 			pks, err := packages.Load(cfg, dir)
// 			if err != nil {
// 				panic(err)
// 			}
// 			return pks[0]
// 		}

// 		v.InitRoot(pkgs, v.Source, packFunc)
// 		v.InitRoot(pkgs, v.Target, packFunc)
// 	}

// 	return mappersMap
// }


// func parseRules(input string) []rule.Any {
// 	var rules []rule.Any

// 	re := regexp.MustCompile(`(\w+)={(.*?)}`)
// 	matches := re.FindAllStringSubmatch(input, -1)

// 	for _, match := range matches {
// 		if len(match) != 3 {
// 			log.Fatalf("invalid rule block: %v", match)
// 		}
// 		ruleType := match[1]
// 		ruleData := match[2]

// 		parser, ok := rule.RuleParsers[ruleType]
// 		if !ok {
// 			log.Fatalf("unknown rule type: %s", ruleType)
// 		}

// 		rule := parser(ruleData)
// 		rules = append(rules, rule)
// 	}

// 	return rules
// }

// func parseQualRule(data string) rule.Any {
// 	resource := regexp.MustCompile(`source="([^"]+)"`)
// 	retarget := regexp.MustCompile(`target="([^"]+)"`)
// 	remname := regexp.MustCompile(`mname="([^"]+)"`)
// 	rempath := regexp.MustCompile(`mpath="([^"]+)"`)

// 	var source, mname, mpath string
// 	smatches := resource.FindStringSubmatch(data)
// 	if len(smatches) == 2 {
// 		source = smatches[1]
// 	}

// 	tmatches := retarget.FindStringSubmatch(data)
// 	if len(tmatches) != 2 {
// 		log.Fatalf("invalid qual format: %s", data)
// 	}

// 	target := tmatches[1]
// 	targetPath := strings.Split(target, ".")
// 	sourcePath := make([]string, len(targetPath))
// 	if strings.HasPrefix(source, ".") {
// 		sourcePath = strings.Split(source, ".")
// 	}

// 	if len(targetPath) != len(sourcePath) {
// 		log.Fatalf("source path must contain the same count of path elements as the target: check rule: @qual={%s}", data)
// 	}

// 	custmnamedata := remname.FindStringSubmatch(data)
// 	if len(custmnamedata) == 2 {
// 		mname = custmnamedata[1]
// 	}

// 	custmnpathdata := rempath.FindStringSubmatch(data)
// 	if len(custmnpathdata) == 2 {
// 		mpath = custmnpathdata[1]
// 	}

// 	return rule.Qual{
// 		SourceName: source,
// 		TargetName: target,
// 		MName:      mname,
// 		MPath:      mpath,
// 	}
// }

// func parseEnumRule(data string) rule.Any {
// 	return rule.Enum{}
// }



// func extractMappingRoot(v *ast.FieldList, imports []struct{ alias, path string }) Struct {
// 	pathByName := func(n string) string {
// 		for _, v := range imports {
// 			if v.alias == n {
// 				return v.path
// 			}
// 		}
// 		return CurrentPath
// 	}

// 	out := Struct{Fields: make(map[string]Field)}
// 	switch expr := v.List[0].Type.(type) {
// 	case *ast.Ident:
// 		out.Path = CurrentPath
// 		out.Name = expr.Name
// 	case *ast.SelectorExpr:
// 		out.Path = pathByName(expr.X.(*ast.Ident).Name)
// 		out.Name = expr.Sel.Name
// 	}

// 	return out
// }

// func extractMapImport(i *ast.ImportSpec) struct{ alias, path string } {
// 	var alias, path string
// 	path = strings.ReplaceAll(i.Path.Value, "\"", "")

// 	if i.Name != nil {
// 		alias = i.Name.Name
// 	} else {
// 		alias = filepath.Base(path)
// 	}

// 	return struct {
// 		alias string
// 		path  string
// 	}{alias: alias, path: path}
// }
