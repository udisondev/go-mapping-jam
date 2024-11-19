//go:generate go run generate.go
package mapper

import (
	"github.com/udisondev/go-mapping-jam/domain"
	d "github.com/udisondev/go-mapping-jam/dto"
)


type Mapper interface {
	// descr
	// sec
    MapPersonToDTO(p domain.Person) d.Person
    MapPersonToDomain(p d.Person) domain.Person
}
