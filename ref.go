package ecs

import "log/slog"

type IsAllocatable[T any] interface {
	Allocate() Ref[T]
	Load(age uint64, id Id) (latestAge uint64, e *T)
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
	Age uint64
	Id  Id
	Ptr *T // make private do not use externally use Get() always
}
type IsEntityPtr[T any] interface {
	IsEntity[T]
	*T
}

func MakeRef[T any, V IsEntityPtr[T]](id Id) Ref[T] {
	return Ref[T]{
		Id: id,
	}
}

func (r *Ref[T]) IsNull() bool {
	return r.Id.IsNull()
}

func (r *Ref[T]) Get() *T {
	if r.Id.IsStored() {
		return r.Ptr
	}

	if !r.Id.IsAllocated() {
		return nil
	}

	loadFn := func(t any) {
		if ai, ok := (t).(IsAllocatable[T]); ok {
			r.Age, r.Ptr = ai.Load(r.Age, r.Id)
		}
	}

	if r.Ptr == nil {
		var t T
		r.Ptr = &t
		r.Age = 0
		loadFn(&t)
	}
	loadFn(r.Ptr)

	return r.Ptr
}

func GetT[
	T any,
	V interface {
		IsAllocatable[T]
		*T
	},
](id Id) (age uint64, e *T) {
	var t T
	return V(&t).Load(0, id)
}

func (r *Ref[T]) Defer() {
	if r.Id.IsNull() {
		return
	}

	if !r.Id.IsAllocated() {
		slog.Warn("Entity already freed.", "id", r.Id)
		return
	}

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
	}(r.Get())

	if !r.Id.IsStored() {
		r.Ptr = nil
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
