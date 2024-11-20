package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"math/rand/v2"
	"regexp"
	"time"

	"strings"

	"golang.org/x/tools/go/packages"

	jen "github.com/dave/jennifer/jen"
)

type FieldName string
type Pack struct {
	Alias string
	Path  string
}
type Direction uint8

type FullType struct {
	ShortPackName string
	StructName    string
}

const (
	Source Direction = iota + 1
	Target
)

type StructField interface {
	isField()
}

func (s *Struct) isField()    {}
func (p *Primetive) isField() {}

type Struct struct {
	Owner    *Struct
	Pack     *Pack
	FullType FullType
	Fields   map[FieldName]StructField
}

type Primetive struct {
	Type string
}

type QualRule struct {
	SourceName       FieldName
	TargetName       FieldName
	CustomMethodName string
}

type EnumRule struct {
	EnumMapping map[string]string
}

type RuleType uint8

const (
	Qual RuleType = iota + 1
	Enum
)

type Rule interface {
	Type() RuleType
}

func (nr QualRule) Type() RuleType { return Qual }
func (er EnumRule) Type() RuleType { return Enum }

// Фабрики для правил
type RuleFactory func(string) Rule

// Глобальная карта парсеров
var ruleParsers = map[string]RuleFactory{}

// Регистрация парсеров
func registerRuleParser(name string, factory RuleFactory) {
	ruleParsers[name] = factory
}

type MapFunc struct {
	Name     string
	Mappable map[Direction]*Struct
	Rules    map[RuleType]Rule
}

const projectName = "github.com/udisondev/go-mapping-jam"

var thisPack = &Pack{Alias: "", Path: projectName + "/mapping"}

var pkgs = make(map[string]*packages.Package)

func main() {
	registerRuleParser("qual", parseQualRule)
	registerRuleParser("enum", parseEnumRule)

	// Путь к файлу с интерфейсом
	interfaceFile := "./mapper/mapper.go"

	// Читаем исходный код
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, interfaceFile, nil, parser.ParseComments)
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

	// mapFuncs := make([]mapper, 0, len(node.Scope.Objects))

	// Ищем интерфейс Mapper и метод

	// Загружаем пакеты проекта
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

	// Генерация кода
	outputFile := "mapper_impl.go"
	generateCodeWithJennifer(outputFile, mapFuncs)

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

func (s Struct) Hash() string {
	return s.Pack.Path + "." + s.FullType.StructName
}

func (p Pack) Dir() string {
	return fmt.Sprintf("./%s", strings.ReplaceAll(string(p.Path), projectName, ""))
}

func (p Pack) Name() string {
	return p.Alias
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
		panic(fmt.Sprintf("type %s not found in package %s", str.FullType.StructName, str.Pack))
	}

	// Проверяем, является ли это именованным типом
	namedType, ok := obj.Type().(*types.Named)
	if !ok {
		panic(fmt.Sprintf("%s is not a named type", str.FullType.StructName))
	}

	// Проверяем, является ли это структурой
	structType, ok := namedType.Underlying().(*types.Struct)
	if !ok {
		panic(fmt.Sprintf("%s is not a struct", str.FullType.StructName))
	}

	// Обрабатываем поля структуры
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldName := FieldName(field.Name())
		str.Fields[fieldName] = m.buildField(str, field)
	}
}

func (m *MapFunc) buildField(owner *Struct, field *types.Var) StructField {
	// Проверяем, является ли поле встроенным типом
	if basicType, ok := field.Type().(*types.Basic); ok {
		return &Primetive{
			Type: basicType.Name(),
		}
	}

	// Проверяем, является ли поле структурой
	if namedType, ok := field.Type().(*types.Named); ok {
		structType, ok := namedType.Underlying().(*types.Struct)
		if !ok {
			// Если это не структура, возвращаем как примитив
			return &Primetive{
				Type: namedType.Obj().Name(),
			}
		}

		// Создаем структуру для вложенного типа
		pack := Pack{Path: namedType.Obj().Pkg().Path()}
		fs := &Struct{
			Owner: owner,
			Pack:  &pack,
			FullType: FullType{
				ShortPackName: namedType.Obj().Pkg().Name(),
				StructName:    namedType.Obj().Name(),
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
	}

	// Если тип неизвестен, возвращаем nil
	panic(fmt.Sprintf("unknown field type: %v", field.Type()))
}

func generateCodeWithJennifer(outputFile string, mapFuncs map[string]MapFunc) {
	f := jen.NewFile("main")
	f.Comment("Code generated by go-mapping-jam. DO NOT EDIT.")

	// Создаём структуру
	f.Type().Id("MapperImpl").Struct()

	// Генерируем функции
	for _, mapFunc := range mapFuncs {
		generateMapFunc(f, mapFunc)
	}

	// Сохраняем файл
	err := f.Save(outputFile)
	if err != nil {
		log.Fatalf("failed to save file: %v", err)
	}
}

func generateMapFunc(f *jen.File, mapFunc MapFunc) {
	sourceStruct := mapFunc.Mappable[Source]
	targetStruct := mapFunc.Mappable[Target]

	// Основной метод
	mapF := f.Func().
		Params(jen.Id("m").Op("*").Id("MapperImpl")).
		Id(mapFunc.Name)
	if sourceStruct.Pack == thisPack {
		mapF.Params(jen.Id("src").Id(sourceStruct.FullType.StructName))
	} else {
		mapF.Params(jen.Id("src").Qual(sourceStruct.Pack.Path, sourceStruct.FullType.StructName))
	}

	if targetStruct.Pack == thisPack {
		mapF.Id(targetStruct.FullType.StructName)
	} else {
		mapF.Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName)
	}

	mapF.BlockFunc(func(g *jen.Group) {
		if targetStruct.Pack == thisPack {
			g.Id("target").Op(":=").Id(targetStruct.FullType.StructName + "{}")
		} else {
			g.Id("target").Op(":=").Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName+"{}")
		}

		for targetFieldName, targetField := range mapFunc.Mappable[Target].Fields {
			qual := mapFunc.Rules[Qual]
			sourceFieldName := targetFieldName
			if q, ok := qual.(QualRule); ok && q.TargetName == targetFieldName {
				sourceFieldName = q.SourceName
			}

			if sourceField, ok := mapFunc.Mappable[Source].Fields[sourceFieldName]; ok {
				// Если это примитивное поле
				if _, ok := targetField.(*Primetive); ok {
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("src").Dot(string(sourceFieldName))
				} else if targetStruct, ok := targetField.(*Struct); ok {
					hash := string(targetStruct.Hash()) + string(sourceField.(*Struct).Hash())
					subMethodName, ok := subMappers[hash]
					if !ok {
						subMethodName = genRandomName(15)
						subMappers[hash] = subMethodName
					}
					// Если это структура
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("m").Dot(subMethodName).Call(jen.Id("src").Dot(string(sourceFieldName)))

					if !ok {
						// Генерируем подметод для вложенных структур
						generateSubMapper(f, subMethodName, sourceField.(*Struct), targetStruct, mapFunc)
					}
				}
			}
		}

		g.Return(jen.Id("target"))
	})
}

var subMappers = make(map[string]string)

func generateSubMapper(f *jen.File, methodName string, sourceStruct *Struct, targetStruct *Struct, mapFunc MapFunc) {
	// Основной метод
	mapF := f.Func().
		Params(jen.Id("m").Op("*").Id("MapperImpl")).
		Id(methodName)
	if sourceStruct.Pack == thisPack {
		mapF.Params(jen.Id("src").Id(sourceStruct.FullType.StructName))
	} else {
		mapF.Params(jen.Id("src").Qual(sourceStruct.Pack.Path, sourceStruct.FullType.StructName))
	}

	if targetStruct.Pack == thisPack {
		mapF.Id(targetStruct.FullType.StructName)
	} else {
		mapF.Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName)
	}

	mapF.BlockFunc(func(g *jen.Group) {
		if targetStruct.Pack == thisPack {
			g.Id("target").Op(":=").Id(targetStruct.FullType.StructName + "{}")
		} else {
			g.Id("target").Op(":=").Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName+"{}")
		}

		for targetFieldName, targetField := range targetStruct.Fields {
			qual := mapFunc.Rules[Qual]
			sourceFieldName := targetFieldName
			if q, ok := qual.(QualRule); ok && q.TargetName == targetFieldName {
				sourceFieldName = q.SourceName
			}

			if sourceField, ok := sourceStruct.Fields[sourceFieldName]; ok {
				if _, ok := sourceField.(*Primetive); ok {
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("src").Dot(string(sourceFieldName))
				} else if nestedSourceStruct, ok := sourceField.(*Struct); ok {
					hash := string(nestedSourceStruct.Hash()) + string(targetField.(*Struct).Hash())
					methodName, ok := subMappers[hash]
					if !ok {
						methodName = genRandomName(15)
						subMappers[hash] = methodName
					}
					g.Id("target").Dot(string(targetFieldName)).Op("=").Id("m").Dot(methodName).Call(jen.Id("src").Dot(string(sourceFieldName)))

					if !ok {
						// Рекурсивно генерируем вложенные подметоды
						generateSubMapper(f, methodName, nestedSourceStruct, targetField.(*Struct), mapFunc)
					}
				}
			}
		}

		g.Return(jen.Id("target"))
	})
}

const charset = "abcdefghijklmnopqrstuvwxyz"

func genRandomName(length int) string {
	seed := time.Now().UnixNano()

	src := rand.NewPCG(uint64(seed), uint64(seed>>32))
	r := rand.New(src)

	// Буфер для сборки строки
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.IntN(len(charset))]
	}
	return string(result)
}

func parseRules(input string) []Rule {
	var rules []Rule

	// Регэксп для поиска блоков типа `qual={...}` или `enum={...}`
	re := regexp.MustCompile(`(\w+)={(.*?)}`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		if len(match) != 3 {
			panic(fmt.Sprintf("invalid rule block: %v", match))
		}
		ruleType := match[1]
		ruleData := match[2]

		// Проверяем, зарегистрирован ли парсер для типа
		parser, ok := ruleParsers[ruleType]
		if !ok {
			panic(fmt.Sprintf("unknown rule type: %s", ruleType))
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
		panic(fmt.Sprintf("invalid qual format: %s", data))
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
		panic(fmt.Sprintf("invalid enum format: %s", data))
	}
	enumMap := make(map[string]string)
	for _, match := range matches {
		enumMap[match[1]] = match[2]
	}
	return EnumRule{EnumMapping: enumMap}
}
