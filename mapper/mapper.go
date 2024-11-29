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

	//@qual={source="Firstname" target=".FirstName"}
	//@qual={target=".LastName" mname="lastNameMapper"}
	//@qual={source="Number" target=".Profile.Phone"}
	MapPersonToDTO(p domain.Person) (d.Person, error)

	//@enum={Simple:Simple Important:Important}
	MapPersonTypeToDto(pt domain.PersonType) d.PersonType

	//@qual={source="FirstName" target=".Firstname"}
	//@qual={source="Phone" target=".Profile.Number"}
	MapPersonToDomain(p d.Person) domain.Person
}

func lastNameMapper(lastName string) string {
	return strings.ToLower(lastName)
}
