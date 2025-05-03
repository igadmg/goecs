package ecs

import (
	"iter"

	"github.com/igadmg/goex/slicesex"
)

type BaseStorage struct {
	id      uint64
	id_pool []Id
	typeId  uint32
	age     uint64

	// Idea: Implement list of retaken ids, fill it if we are iterating over storage (??abstract storage iteration??)
	// if we are iterating further then taken id put it to retaken list
	// iterate retaken list after iterating all entites
	// But seems that makes multithreaded implementation more complex.
	// but anyway we can fallback to command queue anytime
	// Bevys spawn trees looks horrifying tbh

	Ids []Id
}

func MakeBaseStorage(typeId uint32) BaseStorage {
	return BaseStorage{
		typeId: typeId,
	}
}

func (s *BaseStorage) TypeId() uint32 {
	return s.typeId
}

func (s *BaseStorage) Age() uint64 {
	return s.age
}

func (s *BaseStorage) EntityIds() iter.Seq[Id] {
	return func(yield func(id Id) bool) {
		for _, id := range s.Ids {
			if id.IsAllocated() {
				if !yield(id) {
					return
				}
			}
		}

		// Idea: here we can iterate retaken id list to include
		// new object to the processing
	}
}

func (s *BaseStorage) NewId() (id Id) {
	l := len(s.id_pool)
	if l == 0 {
		s.id++
		id = MakeId(s.id, s.typeId)
	} else {
		id = s.id_pool[l-1]
		s.id_pool = s.id_pool[:l-1]
	}

	// init
	id = id.
		Allocate().
		Restore()

	return id
}

func (s *BaseStorage) NewGridId() (id Id) {
	s.id++
	id = MakeId(s.id, s.typeId)

	// init
	id = id.
		Allocate().
		Restore()

	return id
}

func (s *BaseStorage) AllocateId() (age uint64, id Id) {
	id = s.NewId()
	index := int(id.GetId() - 1)

	//glen := min(index, len(s.Ids))
	//grow := index + min(1, glen/5)

	s.age++
	s.Ids = slicesex.Reserve(s.Ids, index+1)
	s.Ids[index] = id

	age = s.age
	return
}

func (s *BaseStorage) AllocateGridIds(size int) (age uint64, start_id Id) {
	si := len(s.Ids)
	s.age++
	s.Ids = slicesex.Reserve(s.Ids, len(s.Ids)+size+1)

	for i := range size {
		id := s.NewGridId()
		s.Ids[si+i] = id
	}

	start_id = s.Ids[si]
	age = s.age
	return
	// return age, start_id

	// return func(yield func(id Id) bool) {
	// 	for i := range size {
	// if !yield(s.Ids[si+i]) {
	// 	return
	// }
	//}
}

func (s *BaseStorage) Free(id Id) Id {
	if !id.IsAllocated() {
		return id
	}

	index := int(id.GetId() - 1)
	if len(s.Ids) <= index {
		return id
	}

	if !s.Ids[index].IsAllocated() {
		return id
	}

	id = id.Free()
	s.age++
	//if !id.IsStored() { // leave stored entities in db
	s.id_pool = append(s.id_pool, id)
	//}
	s.Ids[index] = id
	return id
}

func (s *BaseStorage) EntitiesCount() int64 {
	return int64(s.id) - int64(len(s.id_pool))
}
