package ecs

type BaseEntity struct {
	Id Id
}

//func (e *BaseEntity) Store() {
//	e.Id = e.Id.Store()
//}

//func (e *BaseEntity) Restore() {
//	e.Id = e.Id.Restore()
//}

type FinalBaseEntity[T any] struct {
	BaseEntity
}

func (e FinalBaseEntity[T]) Ref() Ref[T] {
	return Ref[T]{
		Id: e.Id,
	}
}
