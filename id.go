package ecs

import "fmt"

type Id struct {
	value uint64
}

const (
	FlagBits = 2
	TypeBits = 10
	IdBits   = 64 - FlagBits - TypeBits
)

const (
	FlagMaskShift = 64 - FlagBits
	FlagMask      = ((uint64(1) << FlagBits) - 1) << FlagMaskShift
	AllocateBit   = uint64(1) << 63
	StoreBit      = AllocateBit >> 1
	TypeMaskShift = 64 - FlagBits - TypeBits
	TypeMask      = ((uint64(1) << TypeBits) - 1) << TypeMaskShift
	IdMask        = (uint64(1) << IdBits) - 1
)

func MakeId(id uint64, typ uint32) Id {
	return Id{id | (uint64(typ) << TypeMaskShift)}
}

func (id Id) String() string {
	return fmt.Sprintf("Id(%d:%d:%t:%t)", id.GetId(), id.GetType(), id.IsAllocated(), id.IsStored())
}

func (id Id) IsNull() bool {
	return id.value == 0
}

func (id Id) GetId() uint64 {
	return id.value & IdMask
}

func (id Id) SetId(nid uint64) Id {
	return MakeId(nid, id.GetType())
}

func (id Id) GetType() uint32 {
	return (uint32)((id.value & TypeMask) >> TypeMaskShift)
}

func (id Id) IsAllocated() bool {
	return (id.value & AllocateBit) != 0
}

func (id Id) IsStored() bool {
	return (id.value & StoreBit) != 0
}

func (id Id) Allocate() Id {
	id.value = id.value | AllocateBit
	return id
}

func (id Id) Free() Id {
	id.value = id.value & ^AllocateBit
	return id
}

func (id Id) Store() Id {
	id.value = id.value | StoreBit
	return id
}

func (id Id) Restore() Id {
	id.value = id.value & ^StoreBit
	return id
}
