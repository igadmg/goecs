package ecs

import "fmt"

type Id struct {
	value uint64
}

func MakeId(id uint64, typ uint32) Id {
	return Id{id | ((uint64(typ) & 0xffff) << 32)}
}

func (id Id) String() string {
	return fmt.Sprintf("Id(%d:%d:%t:%t)", id.GetId(), id.GetType(), id.IsAllocated(), id.IsStored())
}

func (id Id) IsNull() bool {
	return id.value == 0
}

func (id Id) GetId() uint64 {
	return id.value & 0xffffffff // 32
}

func (id Id) GetType() uint32 {
	return (uint32)((id.value & 0xffff0000000) >> 32) // 32-48
}

const (
	storeBit uint64 = 0x4 << 48
	allocBit uint64 = 0x8 << 48
)

func (id Id) IsAllocated() bool {
	return (id.value & allocBit) != 0
}

func (id Id) IsStored() bool {
	return (id.value & storeBit) != 0
}

func (id Id) Allocate() Id {
	id.value = id.value | allocBit
	return id
}

func (id Id) Free() Id {
	id.value = id.value & ^allocBit
	return id
}

func (id Id) Store() Id {
	id.value = id.value | storeBit
	return id
}

func (id Id) Restore() Id {
	id.value = id.value & ^storeBit
	return id
}
