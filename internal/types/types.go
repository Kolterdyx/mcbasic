package types

import "github.com/Kolterdyx/mcbasic/internal/nbt"

type ValueType interface {
	Primitive() ValueType
	ToString() string
	ToNBT() nbt.Value
	Equals(other ValueType) bool
	GetFieldType(name string) (ValueType, bool)
}
