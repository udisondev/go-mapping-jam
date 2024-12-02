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

	//@qual={source="Firstname" target="t.FirstName"}
	//@qual={target="t.LastName" mname="lastNameMapper"}
	//@qual={source="Number" target="t.Profile.Phone"}
	//@ignore={target="t.Type"}
	MapPersonToDTO(p domain.Person) (d.Person, error)

	//@emapper
	//@enum={Simple:Simple Important:Important}
	MapPersonTypeToDto(pt domain.PersonType) d.PersonType

	//@qual={source="FirstName" target="t.Firstname"}
	//@qual={source="Phone" target="t.Profile.Number"}
	MapPersonToDomain(p d.Person) domain.Person
}

func lastNameMapper(lastName string) string {
	return strings.ToLower(lastName)
}
