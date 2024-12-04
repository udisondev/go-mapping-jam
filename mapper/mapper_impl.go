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

	if src.Firstname != nil {
		target.FirstName = *src.Firstname
	}
	target.LastName = src.LastName
	target.MiddleName = &src.MiddleName
	target.MainAccount = rhlmhkrnwh(src.MainAccount)
	targetAccountSlice := make([]external.Account, 0, len(target.Account))
	for _, it := range src.Account {
		targetAccountSlice = append(targetAccountSlice, rhlmhkrnwh(it))
	}
	target.Account = targetAccountSlice
	if src.Profile != nil {
		target.Profile = maszhjyrtz(*src.Profile)
	}
	target.Projects = src.Projects

	return target
}

func rhlmhkrnwh(src external.Account) external.Account {
	target := external.Account{}

	target.Login = lzguzvqgzh(src.Login)
	target.Password = src.Password

	return target
}

func lzguzvqgzh(src user.Login) user.Login {
	target := user.Login{}

	target.Value = src.Value

	return target
}

func maszhjyrtz(src domain.Profile) dto.Profile {
	target := dto.Profile{}

	target.Phone = src.Number

	return target
}

func MapPersonToDomain(src dto.Person) domain.Person {
	target := domain.Person{}

	target.Firstname = &src.FirstName
	target.LastName = src.LastName
	if src.MiddleName != nil {
		target.MiddleName = *src.MiddleName
	}
	targetAccountSlice := make([]external.Account, 0, len(target.Account))
	for _, it := range src.Account {
		targetAccountSlice = append(targetAccountSlice, rhlmhkrnwh(it))
	}
	target.Account = targetAccountSlice
	target.Projects = src.Projects

	return target
}
