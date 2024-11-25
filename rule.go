package main

type QualRule struct {
	SourceName string
	TargetName string
	MPath      string
	MName      string
}

type CustomMapper struct {
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

type RuleFactory func(string) Rule

var ruleParsers = map[string]RuleFactory{}

func registerRuleParser(name string, factory RuleFactory) {
	ruleParsers[name] = factory
}
