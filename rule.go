package main

type QualRule struct {
	SourceName       string
	TargetName       string
	MethodName string
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

