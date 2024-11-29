package dto

import "github.com/udisondev/go-mapping-jam/external"

type Person struct {
	FirstName string
	LastName string
	MiddleName *string
	Age *int
	Account []external.Account
	Profile Profile
	Type PersonType
	Projects []string
}

type PersonType int

const (
	Simple PersonType = iota + 1
	Important
)

type Profile struct {
	Phone string
}