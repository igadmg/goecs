package ecs

type ArchetypeI[T any] interface {
	Allocate() T
}

type Archetype[T any] struct {
	Id Id
}

func (e Archetype[T]) IsNull() bool {
	return e.Id.IsNull()
}

func (e Archetype[T]) Ref() Ref[T] {
	return Ref[T]{
		Id: e.Id,
	}
}
