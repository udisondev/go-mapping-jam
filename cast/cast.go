package cast

// import (
// 	"fmt"

// 	. "github.com/udisondev/go-mapping-jam/mapp"
// )

// func ToStruct(v any) Struct {
// 	str, ok := v.(Struct)
// 	if !ok {
// 		panic(fmt.Sprintf("'%T' is not a struct", v))
// 	}

// 	return str
// }

// func ToPointerToStruct(v any) Struct {
// 	ptr, ok := v.(Pointer)
// 	if !ok {
// 		panic(fmt.Sprintf("'%T' is not a pointer", v))
// 	}

// 	return ToStruct(ptr.To)
// }

// func ToSliceOfStruct(v any) Struct {
// 	slc, ok := v.(Slice)
// 	if !ok {
// 		panic(fmt.Sprintf("'%T' is not a pointer", v))
// 	}

// 	return ToStruct(slc.Of)
// }
