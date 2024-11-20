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
	target.FirstName = src.FirstName
	target.LastName = src.LastName
	target.MiddleName = src.MiddleName
	target.Age = src.Age
	target.Account = m.oitahzgyrdgjsju(src.Account)
	target.Profile = m.uyoycftrvxwcgnm(src.Profile)
	return target
}
func (m *MapperImpl) oitahzgyrdgjsju(src external.Account) external.Account {
	target := external.Account{}
	target.Login = m.fnbcyfnxhsebldq(src.Login)
	target.Password = src.Password
	return target
}
func (m *MapperImpl) fnbcyfnxhsebldq(src user.Login) user.Login {
	target := user.Login{}
	target.Value = src.Value
	return target
}
func (m *MapperImpl) uyoycftrvxwcgnm(src domain.Profile) dto.Profile {
	target := dto.Profile{}
	target.Phone = src.Phone
	return target
}
func (m *MapperImpl) MapPersonToDomain(src dto.Person) domain.Person {
	target := domain.Person{}
	target.Age = src.Age
	target.Account = m.oitahzgyrdgjsju(src.Account)
	target.Profile = m.qkptczjsxrcwyog(src.Profile)
	target.FirstName = src.FirstName
	target.LastName = src.LastName
	target.MiddleName = src.MiddleName
	return target
}
func (m *MapperImpl) qkptczjsxrcwyog(src dto.Profile) domain.Profile {
	target := domain.Profile{}
	target.Phone = src.Phone
	return target
}
