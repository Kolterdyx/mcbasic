package types

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type StructTypeStruct struct {
	interfaces.ValueType

	Name string
}

func (s StructTypeStruct) Primitive() interfaces.ValueType {
	return s
}
