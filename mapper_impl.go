package main

import (
	domain "github.com/udisondev/go-mapping-jam/domain"
	dto "github.com/udisondev/go-mapping-jam/dto"
	external "github.com/udisondev/go-mapping-jam/external"
	user "github.com/udisondev/go-mapping-jam/user"
)

// Code generated by go-mapping-jam. DO NOT EDIT.
type MapperImpl struct{}

func (m *MapperImpl) MapPersonToDTO(src domain.Person) dto.Person {
	target := dto.Person{}
	target.LastName = src.LastName
	target.MiddleName = src.MiddleName
	target.Age = src.Age
	target.Account = m.litstdrdeovzkfu(src.Account)
	target.Profile = m.yjdlnftqbzglqud(src.Profile)
	target.FirstName = src.Firstname
	return target
}
func (m *MapperImpl) litstdrdeovzkfu(src external.Account) external.Account {
	target := external.Account{}
	target.Login = m.xipeflkohbrubmc(src.Login)
	target.Password = src.Password
	return target
}
func (m *MapperImpl) xipeflkohbrubmc(src user.Login) user.Login {
	target := user.Login{}
	target.Value = src.Value
	return target
}
func (m *MapperImpl) yjdlnftqbzglqud(src domain.Profile) dto.Profile {
	target := dto.Profile{}
	target.Phone = src.Phone
	return target
}
func (m *MapperImpl) MapPersonToDomain(src dto.Person) domain.Person {
	target := domain.Person{}
	target.LastName = src.LastName
	target.MiddleName = src.MiddleName
	target.Age = src.Age
	target.Account = m.litstdrdeovzkfu(src.Account)
	target.Profile = m.jxfwbvcivbowdjg(src.Profile)
	return target
}
func (m *MapperImpl) jxfwbvcivbowdjg(src dto.Profile) domain.Profile {
	target := domain.Profile{}
	target.Phone = src.Phone
	return target
}
