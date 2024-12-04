package gen

import (
	"math/rand/v2"
	"time"

	"github.com/udisondev/go-mapping-jam/mapp"
)

func genRandomName(length int) string {
	seed := time.Now().UnixNano()

	src := rand.NewPCG(uint64(seed), uint64(seed>>32))
	r := rand.New(src)

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.IntN(len(charset))]
	}
	return string(result)

}

func fieldsHash(fs ...mapp.Field) string {
	var hash string
	for _, f := range fs {
		hash += fieldHash(f)
	}

	return hash
}

func fieldHash(f mapp.Field) string {
	return f.Type().Path() + "." + f.Type().TypeName()
}
