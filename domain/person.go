package domain

import "github.com/udisondev/go-mapping-jam/external"

type Person struct {
	Firstname string
	LastName string
	MiddleName string
	Age int64
	Account external.Account
	Profile Profile
}

type Profile struct {
	Phone string
}

