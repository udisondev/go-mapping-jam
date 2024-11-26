package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

var pkgs = make(map[string]*packages.Package)

func parse(filePath string) map[string]Mapper {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("failed to parse file: %v", err)
	}

	mapperImports := make([]struct{ alias, path string }, 0, len(node.Imports))
	for _, v := range node.Imports {
		mapperImports = append(mapperImports, extractMapImport(v))
	}

	for _, v := range extractMethods(node) {
		mappingRules := make(map[RuleType][]Rule)
		if v.Doc != nil {
			for _, mpr := range v.Doc.List {
				rules := parseRules(strings.TrimSpace(strings.ReplaceAll(mpr.Text, "//", "")))
				for _, r := range rules {
					mappingRules[r.Type()] = append(mappingRules[r.Type()], r)
				}
			}
		}

		if fType, ok := v.Type.(*ast.FuncType); ok {
			source := extractMappingRoot(fType.Params, mapperImports)
			target := extractMappingRoot(fType.Results, mapperImports)
			m := Mapper{
				Name:   v.Names[0].Name,
				Source: &source,
				Target: &target,
				Rules:  mappingRules,
			}

			mappersMap[m.Name] = m
		}

	}

	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
		Fset: fset,
	}

	for _, v := range mappersMap {
		packFunc := func(dir string) *packages.Package {
			pks, err := packages.Load(cfg, dir)
			if err != nil {
				panic(err)
			}
			return pks[0]
		}

		v.initRoot(v.Source, packFunc)
		v.initRoot(v.Target, packFunc)
	}

	return mappersMap
}

func (m *Mapper) initRoot(str *Struct, packFunc func(dir string) *packages.Package) {
	pkg, ok := pkgs[str.Path]
	if !ok {
		pkg = packFunc(dirByPath(str.Path))
		pkgs[str.Path] = pkg
	}

	obj := pkg.Types.Scope().Lookup(str.Name)
	if obj == nil {
		log.Fatalf("type %s not found in package %s", str.Name, str.Path)
	}

	namedType, ok := obj.Type().(*types.Named)
	if !ok {
		log.Fatalf("%s is not a named type", str.Name)
	}

	structType, ok := namedType.Underlying().(*types.Struct)
	if !ok {
		log.Fatalf("%s is not a struct", str.Name)
	}

	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldName := field.Name()
		str.Fields[fieldName] = m.buildField(nil, fieldName, field)
	}
}

func dirByPath(p string) string {
	return strings.ReplaceAll(p, projectName, "./")
}

func (m *Mapper) buildField(owner *Field, fieldName string, field *types.Var) *Field {
	switch t := field.Type().(type) {
	case *types.Basic:
		return &Field{
			Owner: owner,
			Name:  fieldName,
			Desc: &Primetive{
				Type: t.Name(),
			}}
	case *types.Named:
		structType, ok := t.Underlying().(*types.Struct)
		if !ok {
			return &Field{
				Owner: owner,
				Name:  fieldName,
				Desc: &Primetive{
					Type: t.Obj().Name(),
				},
			}
		}

		fs := &Field{
			Owner: owner,
			Name:  fieldName,
			Desc: &Struct{
				Path:   t.Obj().Pkg().Path(),
				Name:   t.Obj().Name(),
				Fields: make(map[string]*Field),
			}}

		for i := 0; i < structType.NumFields(); i++ {
			subField := structType.Field(i)
			subFieldName := subField.Name()
			fs.Desc.(*Struct).Fields[subFieldName] = m.buildField(fs, subFieldName, subField)
		}

		return fs

	case *types.Slice:
		b, ok := t.Elem().Underlying().(*types.Basic)
		if ok {
			return &Field{
				Owner: owner,
				Name:  fieldName,
				Desc: &PrimetiveSlice{Primetive: Primetive{
					Type: b.Name(),
				}}}
		}

	case *types.Pointer:
		switch pt := t.Elem().(type) {
		case *types.Basic:
			return &Field{
				Owner: owner,
				Name:  fieldName,
				Desc: &Pointer{
					Ref: &Primetive{
						Type: pt.Name(),
					},
				},
			}
		case *types.Named:
			structType, ok := pt.Underlying().(*types.Struct)
			if !ok {
				return &Field{
					Owner: owner,
					Name:  fieldName,
					Desc: &Pointer{
						Ref: &Primetive{
							Type: pt.Obj().Name(),
						},
					},
				}
			}

			fs := &Field{
				Owner: owner,
				Name:  fieldName,
				Desc: &Pointer{
					Ref: &Struct{
						Path:   pt.Obj().Pkg().Path(),
						Name:   pt.Obj().Name(),
						Fields: make(map[string]*Field),
					}},
			}

			for i := 0; i < structType.NumFields(); i++ {
				subField := structType.Field(i)
				subFieldName := subField.Name()
				fs.Desc.(*Struct).Fields[subFieldName] = m.buildField(fs, subFieldName, subField)
			}

			return fs

		case *types.Slice:
			b, ok := pt.Elem().Underlying().(*types.Basic)
			if ok {
				return &Field{
					Owner: owner,
					Name:  fieldName,
					Desc: &PrimetiveSlice{Primetive: Primetive{
						Type: b.Name(),
					}}}
			}
		}
	}

	log.Fatalf("unknown field type: %v", field.Type())
	return nil
}

func parseRules(input string) []Rule {
	var rules []Rule

	re := regexp.MustCompile(`(\w+)={(.*?)}`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		if len(match) != 3 {
			log.Fatalf("invalid rule block: %v", match)
		}
		ruleType := match[1]
		ruleData := match[2]

		parser, ok := ruleParsers[ruleType]
		if !ok {
			log.Fatalf("unknown rule type: %s", ruleType)
		}

		rule := parser(ruleData)
		rules = append(rules, rule)
	}

	return rules
}

func parseQualRule(data string) Rule {
	resource := regexp.MustCompile(`source="([^"]+)"`)
	retarget := regexp.MustCompile(`target="([^"]+)"`)
	remname := regexp.MustCompile(`mname="([^"]+)"`)
	rempath := regexp.MustCompile(`mpath="([^"]+)"`)

	var source, mname, mpath string
	smatches := resource.FindStringSubmatch(data)
	if len(smatches) == 2 {
		source = smatches[1]
	}

	tmatches := retarget.FindStringSubmatch(data)
	if len(tmatches) != 2 {
		log.Fatalf("invalid qual format: %s", data)
	}

	target := tmatches[1]
	targetPath := strings.Split(target, ".")
	sourcePath := make([]string, len(targetPath))
	if strings.HasPrefix(source, ".") {
		sourcePath = strings.Split(source, ".")
	}

	if len(targetPath) != len(sourcePath) {
		log.Fatalf("source path must contain the same count of path elements as the target: check rule: @qual={%s}", data)
	}

	custmnamedata := remname.FindStringSubmatch(data)
	if len(custmnamedata) == 2 {
		mname = custmnamedata[1]
	}

	custmnpathdata := rempath.FindStringSubmatch(data)
	if len(custmnpathdata) == 2 {
		mpath = custmnpathdata[1]
	}

	return QualRule{
		SourceName: source,
		TargetName: target,
		MName:      mname,
		MPath:      mpath,
	}
}

func parseEnumRule(data string) Rule {
	re := regexp.MustCompile(`(\w+)=([\w]+)`)
	matches := re.FindAllStringSubmatch(data, -1)
	if matches == nil {
		log.Fatalf("invalid enum format: %s", data)
	}
	enumMap := make(map[string]string)
	for _, match := range matches {
		enumMap[match[1]] = match[2]
	}
	return EnumRule{EnumMapping: enumMap}
}

func extractMethods(n ast.Node) []*ast.Field {
	out := []*ast.Field{}
	ast.Inspect(n, func(n ast.Node) bool {
		if iface, ok := n.(*ast.TypeSpec); ok && iface.Name.Name == "Mapper" {
			if ifaceType, ok := iface.Type.(*ast.InterfaceType); ok {
				out = append(out, ifaceType.Methods.List...)
			}
		}
		return true
	})

	return out
}

func extractMappingRoot(v *ast.FieldList, imports []struct{ alias, path string }) Struct {
	pathByName := func(n string) string {
		for _, v := range imports {
			if v.alias == n {
				return v.path
			}
		}
		return currentPath
	}

	out := Struct{Fields: make(map[string]*Field)}
	switch expr := v.List[0].Type.(type) {
	case *ast.Ident:
		out.Path = currentPath
		out.Name = expr.Name
	case *ast.SelectorExpr:
		out.Path = pathByName(expr.X.(*ast.Ident).Name)
		out.Name = expr.Sel.Name
	}

	return out
}

func extractMapImport(i *ast.ImportSpec) struct{ alias, path string } {
	var alias, path string
	path = strings.ReplaceAll(i.Path.Value, "\"", "")

	if i.Name != nil {
		alias = i.Name.Name
	} else {
		alias = filepath.Base(path)
	}

	return struct {
		alias string
		path  string
	}{alias: alias, path: path}
}
