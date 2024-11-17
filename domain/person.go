package domain

type Person struct {
	FirstName string
	LastName string
	MiddleName string
	Age int64
	Account Account
}

type Account struct {
	Login, Password string
}