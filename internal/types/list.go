package types

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type ListTypeStruct struct {
	interfaces.ValueType

	ContentType interfaces.ValueType
}

func NewListType(parent interfaces.ValueType) ListTypeStruct {
	return ListTypeStruct{
		ContentType: parent,
	}
}

func (l ListTypeStruct) Primitive() interfaces.ValueType {
	return l.ContentType.Primitive()
}

func (l ListTypeStruct) ToString() string {
	return l.ContentType.ToString() + "[]"
}
