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
	Projects []string
}

type Profile struct {
	Number string
}

func FirstNameMapper(firstName string) string {
	return strings.ToUpper(firstName)
}