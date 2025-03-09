package ecs

type Entity[T any] struct {
	Id Id
}

func (e Entity[T]) IsNull() bool {
	return e.Id.IsNull()
}

func (e Entity[T]) Ref() Ref[T] {
	return Ref[T]{
		Id: e.Id,
	}
}
 