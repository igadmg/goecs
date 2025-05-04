package ecs

import "fmt"

type Id struct {
	Value uint64
}

type IdType_ struct {
	GetTypeName func(tid uint32) string
}

var IdType IdType_ = IdType_{
	GetTypeName: func(tid uint32) string {
		return "<unk>"
	},
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

var (
	InvalidId Id = Id{0}
)

func MakeId(id uint64, typ uint32) Id {
	return Id{id | (uint64(typ) << TypeMaskShift)}
}

func (id Id) String() string {
	return fmt.Sprintf("Id(%d:%d:%s:%t:%t)", id.GetId(), id.GetType(), IdType.GetTypeName(id.GetType()), id.IsAllocated(), id.IsStored())
}

func (id Id) IsNull() bool {
	return id.Value == 0
}

func (id Id) GetId() uint64 {
	return id.Value & IdMask
}

func (id Id) SetId(nid uint64) Id {
	return MakeId(nid, id.GetType())
}

func (id Id) GetIndex() int {
	return int((id.Value & IdMask) - 1)
}

func (id Id) GetType() uint32 {
	return (uint32)((id.Value & TypeMask) >> TypeMaskShift)
}

func (id Id) IsAllocated() bool {
	return (id.Value & AllocateBit) != 0
}

func (id Id) IsStored() bool {
	return (id.Value & StoreBit) != 0
}

func (id Id) Allocate() Id {
	id.Value = id.Value | AllocateBit
	return id
}

func (id Id) Free() Id {
	id.Value = id.Value & ^AllocateBit
	return id
}

func (id Id) Store() Id {
	id.Value = id.Value | StoreBit
	return id
}

func (id Id) Restore() Id {
	id.Value = id.Value & ^StoreBit
	return id
}
