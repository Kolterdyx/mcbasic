package types

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

type ListTypeStruct struct {
	ValueType

	ContentType ValueType
}

func NewListType(contentType ValueType) ListTypeStruct {
	return ListTypeStruct{
		ContentType: contentType,
	}
}

func (l ListTypeStruct) Primitive() ValueType {
	return l.ContentType.Primitive()
}

func (l ListTypeStruct) ToString() string {
	return l.ContentType.ToString() + "[]"
}

func (l ListTypeStruct) ToNBT() nbt.Value {
	return nbt.NewList()
}

func (l ListTypeStruct) Equals(other ValueType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(ListTypeStruct); ok {
		return l.ContentType.Equals(other.ContentType)
	}
	return false
}
