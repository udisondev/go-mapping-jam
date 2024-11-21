package dto

import "github.com/udisondev/go-mapping-jam/external"

type Person struct {
	FirstName string
	LastName string
	MiddleName string
	Age int64
	Account external.Account
	Profile Profile
	Projects []string
}

type Profile struct {
	Phone string
}