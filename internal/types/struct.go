package types

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type StructTypeStruct struct {
	interfaces.ValueType

	Name string
}

func NewStructType(name string) *StructTypeStruct {
	return &StructTypeStruct{
		Name: name,
	}
}

func (s *StructTypeStruct) Primitive() interfaces.ValueType {
	return s
}

func (s *StructTypeStruct) ToString() string {
	return s.Name
}
