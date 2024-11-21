package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
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

	imports := make(map[string]*Pack)
	initImports := make(map[string]*Pack)
	mapFuncs := make(map[string]MapFunc)
	for _, v := range node.Imports {
		imp := extractMapImport(v)
		imports[imp.Path] = &imp
		initImports[imp.Alias] = &imp
	}

	astMethods := extractMethods(node)

	for _, v := range astMethods {
		mappingRules := make(map[RuleType]Rule)
		if v.Doc != nil {
			for _, mpr := range v.Doc.List {
				rules := parseRules(strings.TrimSpace(strings.ReplaceAll(mpr.Text, "//", "")))
				for _, r := range rules {
					mappingRules[r.Type()] = r
				}
			}
		}

		if fType, ok := v.Type.(*ast.FuncType); ok {
			source := extractMappingRoot(fType.Params, initImports)
			target := extractMappingRoot(fType.Results, initImports)
			m := MapFunc{
				Name:     v.Names[0].Name,
				Mappable: map[Direction]*Struct{Source: &source, Target: &target},
				Rules:    mappingRules,
			}

			mapFuncs[m.Name] = m
		}
	}

	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports | packages.NeedSyntax,
		Fset: fset,
	}

	for _, v := range mapFuncs {
		packFunc := func(dir string) *packages.Package {
			pks, err := packages.Load(cfg, dir)
			if err != nil {
				panic(err)
			}
			return pks[0]
		}

		v.initRoot(v.Mappable[Source], packFunc)
		v.initRoot(v.Mappable[Target], packFunc)
	}

	return mapFuncs
}


func (m *MapFunc) initRoot(str *Struct, packFunc func(dir string) *packages.Package) {
	// Загружаем пакет
	pkg, ok := pkgs[str.Pack.Path]
	if !ok {
		pkg = packFunc(string(str.Pack.Dir()))
		pkgs[str.Pack.Path] = pkg
	}

	// Находим тип в пакете
	obj := pkg.Types.Scope().Lookup(str.FullType.StructName)
	if obj == nil {
		log.Fatalf("type %s not found in package %s", str.FullType.StructName, str.Pack)
	}

	// Проверяем, является ли это именованным типом
	namedType, ok := obj.Type().(*types.Named)
	if !ok {
		log.Fatalf("%s is not a named type", str.FullType.StructName)
	}

	// Проверяем, является ли это структурой
	structType, ok := namedType.Underlying().(*types.Struct)
	if !ok {
		log.Fatalf("%s is not a struct", str.FullType.StructName)
	}

	// Обрабатываем поля структуры
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldName := FieldName(field.Name())
		str.Fields[fieldName] = m.buildField(str, field)
	}
}

func (m *MapFunc) buildField(owner *Struct, field *types.Var) StructField {

	switch t := field.Type().(type) {
	case *types.Basic:
		return &Primetive{
			Type: t.Name(),
		}
	case *types.Named: 
		structType, ok := t.Underlying().(*types.Struct)
		if !ok {
			// Если это не структура, возвращаем как примитив
			return &Primetive{
				Type: t.Obj().Name(),
			}
		}

		// Создаем структуру для вложенного типа
		pack := Pack{Path: t.Obj().Pkg().Path()}
		fs := &Struct{
			Owner: owner,
			Pack:  &pack,
			FullType: FullType{
				ShortPackName: t.Obj().Pkg().Name(),
				StructName:    t.Obj().Name(),
			},
			Fields: make(map[FieldName]StructField),
		}

		// Рекурсивно обрабатываем поля структуры
		for i := 0; i < structType.NumFields(); i++ {
			subField := structType.Field(i)
			subFieldName := FieldName(subField.Name())
			fs.Fields[subFieldName] = m.buildField(fs, subField)
		}

		return fs

	case *types.Slice: 
		b, ok := t.Elem().Underlying().(*types.Basic)
		if ok {
			return &PrimetiveSlice{Type: b.Name()}
		}
	}

	// Если тип неизвестен, возвращаем nil
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
		SourceName: FieldName(matches[1]),
		TargetName: FieldName(matches[2]),
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

func extractMappingRoot(v *ast.FieldList, impMap map[string]*Pack) Struct {
	out := Struct{Fields: make(map[FieldName]StructField)}
	switch expr := v.List[0].Type.(type) {
	case *ast.Ident:
		out.FullType = FullType{StructName: expr.Name}
		out.Pack = thisPack
	case *ast.SelectorExpr:
		pack := impMap[expr.X.(*ast.Ident).Name]
		out.FullType = FullType{ShortPackName: expr.X.(*ast.Ident).Name, StructName: expr.Sel.Name}
		out.Pack = pack
	}

	return out
}


func extractMapImport(i *ast.ImportSpec) Pack {
	out := Pack{}
	path := strings.ReplaceAll(i.Path.Value, "\"", "")
	out.Path = path

	if i.Name != nil {
		out.Alias = i.Name.Name
		return out
	}

	pathElements := strings.Split(path, "/")
	lastPathElement := pathElements[len(pathElements)-1]
	out.Alias = lastPathElement

	return out
}