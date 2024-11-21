package main

import (
	"fmt"
	"strings"
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


func (s Struct) Hash() string {
	return s.Pack.Path + "." + s.FullType.StructName
}

func (p Pack) Dir() string {
	return fmt.Sprintf("./%s", strings.ReplaceAll(string(p.Path), projectName, ""))
}

func (p Pack) Name() string {
	return p.Alias
}

