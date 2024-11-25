//go:generate go run generate.go
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
// TODO убрать impl call with package
// TODO cast


// TODO добавить source path
// TODO expr
type Mapper interface {

	//@qual={source="Firstname" target=".FirstName" mname="firstNameMapper"}
	//@qual={source="Number" target=".Profile.Phone"}
	//@err(source="Firstname" target=".FirstName" errf="")
	MapPersonToDTO(p domain.Person) (d.Person, error)

	//@qual={source="Firstname" target=".FirstName"}
	//@qual={source=".Prof.Phone" target=".Profile.Number"}
	//@qual={source="Prof" target=".Profile"}
	MapPersonToDomain(p d.Person) domain.Person
}

func firstNameMapper(firstName string) string {
	return strings.ToUpper(firstName)
}

