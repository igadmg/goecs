package ecs

import "log/slog"

type IsAllocatable[T any] interface {
	Allocate() Ref[T]
}

type IsLoadable[T any] interface {
	Load(age uint64, id Id) (latestAge uint64, e T)
}

type IsFreeable interface {
	Free()
}

type IsStorable interface {
	Store()
	Restore()
}

type IsDeferable interface {
	Defer()
}

type IsEntity[T any] interface {
	IsAllocatable[T]
	IsStorable
	IsDeferable
}

// Ref represents reference to ECS entity.
// It stores entity Id and Age of last update.
// Internal pointer Ptr sotres pointer to last updated value.
// Ref can be Defered when no longer in need.
// Defer will free entity from ECS world, but will keep last
type Ref[T any] struct {
	Age    uint64
	Id     Id
	Ptr    T // make private do not use externally use Get() always
	isNull bool
}
type IsEntityPtr[T any] interface {
	IsEntity[T]
	*T
}

func MakeRef[T any, V IsEntityPtr[T]](id Id) Ref[T] {
	return Ref[T]{
		Id:     id,
		isNull: true,
	}
}

func (r Ref[T]) IsNull() bool {
	return r.Id.IsNull() || r.isNull
}

func (r *Ref[T]) Get() T {
	if r.Id.IsStored() {
		return r.Ptr
	}

	if !r.Id.IsAllocated() {
		var t T
		return t
	}

	if r.isNull {
		r.Age = 0
	}

	var t any = r.Ptr
	if ai, ok := (t).(IsLoadable[T]); ok {
		r.Age, r.Ptr = ai.Load(r.Age, r.Id)
	}

	r.isNull = false

	return r.Ptr
}

func GetT[T IsLoadable[T]](id Id) (age uint64, e T) {
	var t T
	return t.Load(0, id)
}

func (r *Ref[T]) Defer() {
	if r.Id.IsNull() {
		return
	}

	if !r.Id.IsAllocated() {
		slog.Warn("Entity already freed.", "id", r.Id)
		return
	}

	r.Get()
	func(t any) {
		if t == nil {
			slog.Warn("Entity already freed.", "id", r.Id)
			return
		}

		if di, ok := t.(IsDeferable); ok {
			di.Defer()
		}
		if fi, ok := t.(IsFreeable); ok {
			fi.Free()
		}
	}(&r.Ptr)

	if !r.Id.IsStored() {
		r.isNull = true
	}
	r.Id = r.Id.Free()
}

func (r *Ref[T]) Store() {
	if r.Id.IsStored() {
		return
	}

	func(t any) {
		if si, ok := t.(IsStorable); ok {
			si.Store()
		}
	}(r.Get()) // this Loads Ptr which is saved.

	r.Id = r.Id.Store()
	r.Defer()
}

func (r *Ref[T]) Restore() {
	if !r.Id.IsStored() {
		return
	}

	func(t any) {
		if ai, ok := t.(IsAllocatable[T]); ok {
			*r = ai.Allocate()
		}
		if si, ok := t.(IsStorable); ok {
			si.Restore()
		}
	}(r.Ptr)

	r.Id = r.Id.Restore()
}
