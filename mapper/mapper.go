//go:generate go run generate.go
package mapper

import (
	"strings"

	"github.com/udisondev/go-mapping-jam/domain"
	d "github.com/udisondev/go-mapping-jam/dto"
)


type Mapper interface {
	//@qual={source="Firstname" target=".FirstName" mname="firstNameMapper"}
	//@qual={source="Number" target=".Profile.Phone"}
    MapPersonToDTO(p domain.Person) d.Person

	//@qual={source="Firstname" target=".FirstName"}
	//@qual={source="Phone" target=".Profile.Number"}
    MapPersonToDomain(p d.Person) domain.Person
}

func firstNameMapper(firstName string) string {
	return strings.ToUpper(firstName)
}