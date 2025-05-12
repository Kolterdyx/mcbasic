package types

import "github.com/Kolterdyx/mcbasic/internal/nbt"

type NamedType struct {
	ValueType
	name string
}

func NewNamedType(name string) ValueType {
	return &NamedType{
		name: name,
	}
}

func (n *NamedType) Primitive() ValueType {
	return n
}

func (n *NamedType) ToString() string {
	return n.name
}

func (n *NamedType) ToNBT() nbt.Value {
	return nbt.NewString(n.name)
}

func (n *NamedType) Equals(other ValueType) bool {
	if otherType, ok := other.(*NamedType); ok {
		return n.name == otherType.name
	}
	return false
}
