package ecs

import (
	"encoding/gob"
	"iter"

	"github.com/igadmg/goex/slicesex"
)

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

func Storage[T any](id Id) (T, bool) {
	i := id.GetType()
	if int(i) < len(registeredTypes) {
		t, ok := registeredTypes[i].(T)
		return t, ok
	}

	var t T
	return t, false
}

func StorageSeq[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, rt := range registeredTypes {
			if t, ok := rt.(T); ok {
				if !yield(t) {
					return
				}
			}
		}
	}
}

type WithAny interface {
	WithAny(id Id, do func(any))
}

func Interface[T any](id Id) (r T, rok bool) {
	if s, ok := Storage[WithAny](id); ok {
		s.WithAny(id, func(i any) {
			if t, ok := i.(T); ok {
				r = t
				rok = true
			}
		})
	}

	return
}

func RefId[T any](u Ref[T]) Id { return u.Id }

type Packer interface {
	PrePack()
	Pack()
}

type Saver interface {
	Save(w *gob.Encoder) error
}

type Loader interface {
	Load(w *gob.Decoder) error
}
