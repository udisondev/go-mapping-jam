package mapper

import (
	"strings"

	"github.com/udisondev/go-mapping-jam/domain"
	d "github.com/udisondev/go-mapping-jam/dto"
)

// TODO добавить обработку slice
// TODO добавить обработку map
// TODO добавить обработку enum (with default)
// TODO work with err

// TODO добавить source path
// TODO expr
type Mapper interface {

	//@qual -s=Firstname -t=t.FirstName
	//@qual -t=t.LastName -mn=lastNameMapper
	//@qual -s=Number -t=t.Profile.Phone
	//@ignore -t=t.Type
	MapPersonToDTO(p domain.Person) (d.Person, error)

	//@emapper
	//@enum Simple=Simple Important=Important
	MapPersonTypeToDto(pt domain.PersonType) d.PersonType

	//@qual -s=FirstName -t=t.Firstname
	//@qual -s=Phone -t=t.Profile.Number
	MapPersonToDomain(p d.Person) domain.Person
}

func lastNameMapper(lastName string) string {
	return strings.ToLower(lastName)
}
