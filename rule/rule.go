//go:generate go-enum
package rule

type Qual struct {
	SourceName string
	TargetName string
	MPath      string
	MName      string
}

type CustomMapper struct {
}

type Enum struct {
	Source EnumItem
	Target EnumItem
	Map    map[string]string
}

type EnumItem struct {
	Name string
	Path string
}

// Type ENUM(qual, enum)
type Type uint8

type Any interface {
	Type() Type
}

func (nr Qual) Type() Type     { return TypeQual }
func (er Enum) Type() Type { return TypeEnum }

type RuleFactory func(string) Any

var RuleParsers = map[string]RuleFactory{}

func RegisterRuleParser(name string, factory RuleFactory) {
	RuleParsers[name] = factory
}
