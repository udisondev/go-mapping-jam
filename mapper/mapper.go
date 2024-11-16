//go:generate go run generate.go
package mapper

import (
	"github.com/udisondev/go-mapping-jam/domain"
	d "github.com/udisondev/go-mapping-jam/dto"
)


type Mapper interface {
	// descr
	// sec
    MapPerson(p domain.Person) d.Person
}
