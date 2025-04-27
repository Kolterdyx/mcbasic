package types

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"reflect"
)

type ListTypeStruct struct {
	interfaces.ValueType

	Parent interfaces.ValueType
}

func (l ListTypeStruct) IsType(other interfaces.ValueType) bool {
	return reflect.TypeOf(l) == reflect.TypeOf(other)
}

func (l ListTypeStruct) Primitive() interfaces.ValueType {
	return l.Parent.Primitive()
}

func (l ListTypeStruct) ToString() string {
	return l.Parent.ToString() + "[]"
}
