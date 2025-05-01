package types

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

type StructTypeStruct struct {
	interfaces.ValueType

	Name   string
	fields map[string]interfaces.ValueType
}

func NewStructType(name string) StructTypeStruct {
	return StructTypeStruct{
		Name:   name,
		fields: make(map[string]interfaces.ValueType),
	}
}

func (s StructTypeStruct) SetField(name string, value interfaces.ValueType) {
	if s.fields == nil {
		s.fields = make(map[string]interfaces.ValueType)
	}
	s.fields[name] = value
}

func (s StructTypeStruct) GetField(name string) (interfaces.ValueType, bool) {
	value, ok := s.fields[name]
	return value, ok
}

func (s StructTypeStruct) Primitive() interfaces.ValueType {
	return s
}

func (s StructTypeStruct) ToString() string {
	return s.Name
}

func (s StructTypeStruct) ToNBT() nbt.Value {
	compound := nbt.NewCompound()
	for name, value := range s.fields {
		compound.Set(name, value.ToNBT())
	}
	return compound
}

func (s StructTypeStruct) Size() int {
	return len(s.fields)
}
