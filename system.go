package ecs

type System interface {
	IsDeferable

	Update(dt float32)
}
