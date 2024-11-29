// Code generated by go-enum DO NOT EDIT.
// Version:
// Revision:
// Build Date:
// Built By:

package main

import (
	"errors"
	"fmt"
)

const (
	// RuleTypeQual is a RuleType of type Qual.
	RuleTypeQual RuleType = iota
	// RuleTypeEnum is a RuleType of type Enum.
	RuleTypeEnum
)

var ErrInvalidRuleType = errors.New("not a valid RuleType")

const _RuleTypeName = "qualenum"

var _RuleTypeMap = map[RuleType]string{
	RuleTypeQual: _RuleTypeName[0:4],
	RuleTypeEnum: _RuleTypeName[4:8],
}

// String implements the Stringer interface.
func (x RuleType) String() string {
	if str, ok := _RuleTypeMap[x]; ok {
		return str
	}
	return fmt.Sprintf("RuleType(%d)", x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x RuleType) IsValid() bool {
	_, ok := _RuleTypeMap[x]
	return ok
}

var _RuleTypeValue = map[string]RuleType{
	_RuleTypeName[0:4]: RuleTypeQual,
	_RuleTypeName[4:8]: RuleTypeEnum,
}

// ParseRuleType attempts to convert a string to a RuleType.
func ParseRuleType(name string) (RuleType, error) {
	if x, ok := _RuleTypeValue[name]; ok {
		return x, nil
	}
	return RuleType(0), fmt.Errorf("%s is %w", name, ErrInvalidRuleType)
}