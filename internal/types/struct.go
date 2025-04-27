package types

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"reflect"
)

type StructTypeStruct struct {
	interfaces.ValueType

	Name string
}

func (s StructTypeStruct) Primitive() interfaces.ValueType {
	return s
}

func (s StructTypeStruct) IsType(other interfaces.ValueType) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(other)
}

func (s StructTypeStruct) ToString() string {
	return s.Name
}
