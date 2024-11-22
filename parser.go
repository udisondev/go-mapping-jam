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

func parse(filePath string) map[string]MapFunc {
	// Читаем исходный код
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
			m := MapFunc{
				Name:   v.Names[0].Name,
				Source: &source,
				Target: &target,
				Rules: mappingRules,
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

func (m *MapFunc) initRoot(str *Struct, packFunc func(dir string) *packages.Package) {
	// Загружаем пакет
	pkg, ok := pkgs[str.Path]
	if !ok {
		pkg = packFunc(dirByPath(str.Path))
		pkgs[str.Path] = pkg
	}

	// Находим тип в пакете
	obj := pkg.Types.Scope().Lookup(str.Name)
	if obj == nil {
		log.Fatalf("type %s not found in package %s", str.Name, str.Path)
	}

	// Проверяем, является ли это именованным типом
	namedType, ok := obj.Type().(*types.Named)
	if !ok {
		log.Fatalf("%s is not a named type", str.Name)
	}

	// Проверяем, является ли это структурой
	structType, ok := namedType.Underlying().(*types.Struct)
	if !ok {
		log.Fatalf("%s is not a struct", str.Name)
	}

	// Обрабатываем поля структуры
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldName := field.Name()
		str.Fields[fieldName] = m.buildField(nil, fieldName, field)
	}
}

func dirByPath(p string) string {
	return strings.ReplaceAll(p, projectName, "./")
}

func (m *MapFunc) buildField(owner *Field, fieldName string, field *types.Var) *Field {
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
			// Если это не структура, возвращаем как примитив
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

		// Рекурсивно обрабатываем поля структуры
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
	}

	log.Fatalf("unknown field type: %v", field.Type())
	return nil
}

func parseRules(input string) []Rule {
	var rules []Rule

	// Регэксп для поиска блоков типа `qual={...}` или `enum={...}`
	re := regexp.MustCompile(`(\w+)={(.*?)}`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		if len(match) != 3 {
			log.Fatalf("invalid rule block: %v", match)
		}
		ruleType := match[1]
		ruleData := match[2]

		// Проверяем, зарегистрирован ли парсер для типа
		parser, ok := ruleParsers[ruleType]
		if !ok {
			log.Fatalf("unknown rule type: %s", ruleType)
		}

		// Парсим правило
		rule := parser(ruleData)
		rules = append(rules, rule)
	}

	return rules
}

// Парсер для QualRule
func parseQualRule(data string) Rule {
	// Используем регэксп для извлечения данных
	re := regexp.MustCompile(`source="([^"]+)"\s+target="([^"]+)"`)
	
	matches := re.FindStringSubmatch(data)
	if len(matches) != 3 {
		log.Fatalf("invalid qual format: %s", data)
	}
	return QualRule{
		SourceName: matches[1],
		TargetName: matches[2],
	}
}

// Парсер для EnumRule
func parseEnumRule(data string) Rule {
	// Используем регэксп для извлечения пар ключ=значение
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
