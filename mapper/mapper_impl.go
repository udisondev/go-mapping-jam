package mapper

import (
	domain "github.com/udisondev/go-mapping-jam/domain"
	dto "github.com/udisondev/go-mapping-jam/dto"
	external "github.com/udisondev/go-mapping-jam/external"
	user "github.com/udisondev/go-mapping-jam/user"
)

// Code generated by go-mapping-jam. DO NOT EDIT.
func MapPersonToDomain(src dto.Person) domain.Person {
	target := domain.Person{}
	target.Firstname = src.FirstName
	target.LastName = src.LastName
	target.MiddleName = src.MiddleName
	target.Age = int64(src.Age)
	target.Account = lvgiqixfbgocxbt(src.Account)
	target.Profile = nxfhbypcezzcwif(src.Profile)
	return target
}
func lvgiqixfbgocxbt(src external.Account) external.Account {
	target := external.Account{}
	target.Password = src.Password
	target.Login = faffxhgmeztupwy(src.Login)
	return target
}
func faffxhgmeztupwy(src user.Login) user.Login {
	target := user.Login{}
	target.Value = src.Value
	return target
}
func nxfhbypcezzcwif(src dto.Profile) domain.Profile {
	target := domain.Profile{}
	target.Number = src.Phone
	return target
}
func MapPersonToDTO(src domain.Person) dto.Person {
	target := dto.Person{}
	target.FirstName = domain.FirstNameMapper(src.Firstname)
	target.LastName = lastNameMapper(src.LastName)
	target.MiddleName = src.MiddleName
	target.Age = int(src.Age)
	target.Account = lvgiqixfbgocxbt(src.Account)
	target.Profile = podbsfxybzflmzq(src.Profile)
	return target
}
func podbsfxybzflmzq(src domain.Profile) dto.Profile {
	target := dto.Profile{}
	target.Phone = src.Number
	return target
}
