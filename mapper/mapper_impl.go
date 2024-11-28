package mapper

import (
	domain "github.com/udisondev/go-mapping-jam/domain"
	dto "github.com/udisondev/go-mapping-jam/dto"
	external "github.com/udisondev/go-mapping-jam/external"
	user "github.com/udisondev/go-mapping-jam/user"
)

// Code generated by go-mapping-jam. DO NOT EDIT.
func MapPersonToDTO(src domain.Person) dto.Person {
	target := dto.Person{}
	if src.Profile != nil {
		target.Profile = wxbnfqbbyyjsqdq(*src.Profile)
	}
	target.Projects = src.Projects
	if src.Firstname != nil {
		target.FirstName = *src.Firstname
	}
	target.LastName = lastNameMapper(src.LastName)
	target.MiddleName = &src.MiddleName
	target.Age = src.Age
	if src.Account != nil {
		tnsaylfivpjxdfkResult := tnsaylfivpjxdfk(*src.Account)
		target.Account = &tnsaylfivpjxdfkResult
	}
	return target
}
func wxbnfqbbyyjsqdq(src domain.Profile) dto.Profile {
	target := dto.Profile{}
	target.Phone = src.Number
	return target
}
func tnsaylfivpjxdfk(src external.Account) external.Account {
	target := external.Account{}
	target.Login = vwqfztcevpdmvfz(src.Login)
	target.Password = src.Password
	return target
}
func vwqfztcevpdmvfz(src user.Login) user.Login {
	target := user.Login{}
	target.Value = src.Value
	return target
}
func MapPersonToDomain(src dto.Person) domain.Person {
	target := domain.Person{}
	target.Firstname = &src.FirstName
	target.LastName = src.LastName
	if src.MiddleName != nil {
		target.MiddleName = *src.MiddleName
	}
	target.Age = src.Age
	if src.Account != nil {
		tnsaylfivpjxdfkResult := tnsaylfivpjxdfk(*src.Account)
		target.Account = &tnsaylfivpjxdfkResult
	}
	saztzqdgjfioecyResult := saztzqdgjfioecy(src.Profile)
	target.Profile = &saztzqdgjfioecyResult
	target.Projects = src.Projects
	return target
}
func saztzqdgjfioecy(src dto.Profile) domain.Profile {
	target := domain.Profile{}
	target.Number = src.Phone
	return target
}
