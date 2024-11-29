//go:generate go-enum
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
	Source EnumItem
	Target EnumItem
	Map map[string]string
}

type EnumItem struct {
	Name string
	Path string
}

// RuleType ENUM(qual, enum)
type RuleType uint8

type Rule interface {
	Type() RuleType
}

func (nr QualRule) Type() RuleType { return RuleTypeQual }
func (er EnumRule) Type() RuleType { return RuleTypeEnum }

type RuleFactory func(string) Rule

var ruleParsers = map[string]RuleFactory{}

func registerRuleParser(name string, factory RuleFactory) {
	ruleParsers[name] = factory
}
