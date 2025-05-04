package ecs

import "github.com/igadmg/goex/slicesex"

type MetaTag struct{} // Used to add tags to structs

//type Query[T any] struct {
//	Age   uint64
//	Items []T
//}

var registeredTypes = []any{}

func RequestRegisterTypes(num int) []any {
	registeredTypes = slicesex.Reserve(registeredTypes, num)
	return registeredTypes
}

func Static[T any](id Id) T {
	i := id.GetType()
	if int(i) < len(registeredTypes) {
		return registeredTypes[i].(T)
	}

	var t T
	return t
}
