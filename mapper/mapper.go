package mapper

import (
	"strings"

	"github.com/udisondev/go-mapping-jam/domain"
	d "github.com/udisondev/go-mapping-jam/dto"
)

// TODO добавить обработку slice
// TODO добавить обработку map
// TODO добавить обработку pointer
// TODO добавить обработку enum (with default)
// TODO work with err
// TODO cast

// TODO добавить source path
// TODO expr
type Mapper interface {

	//@qual={source="Firstname" target=".FirstName" mname="FirstNameMapper" mpath="github.com/udisondev/go-mapping-jam/domain"}
	//@qual={target=".LastName" mname="lastNameMapper"}
	//@qual={target=".Age" mname="int"}
	//@qual={source="Number" target=".Profile.Phone"}
	MapPersonToDTO(p domain.Person) (d.Person, error)

	//@qual={source="FirstName" target=".Firstname"}
	//@qual={source="Phone" target=".Profile.Number"}
	//@qual={target=".Age" mname="int64"}
	MapPersonToDomain(p d.Person) domain.Person
}

func lastNameMapper(lastName string) string {
	return strings.ToLower(lastName)
}
