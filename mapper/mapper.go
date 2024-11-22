//go:generate go run generate.go
package mapper

import (
	"github.com/udisondev/go-mapping-jam/domain"
	d "github.com/udisondev/go-mapping-jam/dto"
)


type Mapper interface {
	//@qual={source="Firstname" target=".FirstName"}
	//@qual={source="Number" target=".Profile.Phone"}
    MapPersonToDTO(p domain.Person) d.Person

	//@qual={source="Firstname" target=".FirstName"}
	//@qual={source="Number" target=".Profile.Phone"}
    MapPersonToDomain(p d.Person) domain.Person
}