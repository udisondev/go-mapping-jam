package domain

import (
	"strings"

	"github.com/udisondev/go-mapping-jam/external"
)

type Person struct {
	Firstname *string
	LastName string
	MiddleName string
	Age *int
	Account []external.Account
	Profile *Profile
	Type PersonType
	Projects []string
}

type PersonType uint8

const (
	Simple PersonType = iota + 1
	Important
)

type Profile struct {
	Number string
}

func FirstNameMapper(firstName string) string {
	return strings.ToUpper(firstName)
}