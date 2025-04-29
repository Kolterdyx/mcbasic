package types

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type ListTypeStruct struct {
	interfaces.ValueType

	Parent interfaces.ValueType
}

func (l ListTypeStruct) Primitive() interfaces.ValueType {
	return l.Parent.Primitive()
}

func (l ListTypeStruct) ToString() string {
	return l.Parent.ToString() + "[]"
}
