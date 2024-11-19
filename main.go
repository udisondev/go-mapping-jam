package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"

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

type Rule struct {
	Value string
}

type MapFunc struct {
	Name     string
	Mappable map[Direction]*Struct
	Rules    []Rule
}

const projectName = "github.com/udisondev/go-mapping-jam"

var thisPack = &Pack{Alias: "", Path: projectName + "/mapping"}

var pkgs = make(map[string]*packages.Package)

func main() {
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
		var mappingRules []Rule
		if v.Doc != nil {
			for _, mpr := range v.Doc.List {
				mappingRules = append(mappingRules, Rule{
					Value: mpr.Text,
				})
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
			g.Id("target").Op(":=").Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName + "{}")
		}

		for sourceFieldName, sourceField := range mapFunc.Mappable[Source].Fields {
			if targetField, ok := mapFunc.Mappable[Target].Fields[sourceFieldName]; ok {
				// Если это примитивное поле
				if _, ok := sourceField.(*Primetive); ok {
					g.Id("target").Dot(string(sourceFieldName)).Op("=").Id("src").Dot(string(sourceFieldName))
				} else if sourceStruct, ok := sourceField.(*Struct); ok {
					// Если это структура
					subMethodName := fmt.Sprintf("Map%sTo%s", sourceStruct.FullType.StructName, targetField.(*Struct).FullType.StructName)
					g.Id("target").Dot(string(sourceFieldName)).Op("=").Id("m").Dot(subMethodName).Call(jen.Id("src").Dot(string(sourceFieldName)))

					// Генерируем подметод для вложенных структур
					generateSubMapper(f, subMethodName, sourceStruct, targetField.(*Struct), mapFunc)
				}
			}
		}

		g.Return(jen.Id("target"))
	})
}

var subMappers = make(map[string]struct{})

func generateSubMapper(f *jen.File, methodName string, sourceStruct *Struct, targetStruct *Struct, mapFunc MapFunc) {
	hash := string(sourceStruct.Hash()) + string(targetStruct.Hash())
	if _, ok := subMappers[hash]; ok {
		return
	}
	subMappers[hash] = struct{}{}
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
			g.Id("target").Op(":=").Qual(targetStruct.Pack.Path, targetStruct.FullType.StructName + "{}")
		}

		for sourceFieldName, sourceField := range sourceStruct.Fields {
			if targetField, ok := targetStruct.Fields[sourceFieldName]; ok {
				if _, ok := sourceField.(*Primetive); ok {
					g.Id("target").Dot(string(sourceFieldName)).Op("=").Id("src").Dot(string(sourceFieldName))
				} else if nestedSourceStruct, ok := sourceField.(*Struct); ok {
					nestedMethodName := fmt.Sprintf("Map%sTo%s", nestedSourceStruct.FullType.StructName, targetField.(*Struct).FullType.StructName)
					g.Id("target").Dot(string(sourceFieldName)).Op("=").Id("m").Dot(nestedMethodName).Call(jen.Id("src").Dot(string(sourceFieldName)))

					// Рекурсивно генерируем вложенные подметоды
					generateSubMapper(f, nestedMethodName, nestedSourceStruct, targetField.(*Struct), mapFunc)
				}
			}
		}

		g.Return(jen.Id("target"))
	})
}