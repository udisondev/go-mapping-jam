//go:generate go-enum
package mapp

import (
	"fmt"
	"strings"
)

// RuleType ENUM(qual, ignore)
type RuleType uint8

// RuleArg ENUM(source, target, mname, mpath)
type RuleArg uint8

type Rule struct {
	spec string
	args map[RuleArg]string
}

func (r Rule) Type() RuleType {
	if strings.HasPrefix(r.spec, "@qual") {
		return RuleTypeQual
	}
	if strings.HasPrefix(r.spec, "@ignore") {
		return RuleTypeIgnore
	}

	panic(fmt.Sprintf("can't define rule type for rule: %s", r.spec))
}

func (r Rule) Value() string {
	return r.spec
}

func (r Rule) Arg(arg  RuleArg) (string, bool) {
	val, ok := r.args[arg]
	if ok {
		return val, ok
	}

	return "", false
}
