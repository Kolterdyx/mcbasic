package types

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

type StructTypeStruct struct {
	ValueType

	Name   string
	fields map[string]ValueType
}

func NewStructType(name string) StructTypeStruct {
	return StructTypeStruct{
		Name:   name,
		fields: make(map[string]ValueType),
	}
}

func (s StructTypeStruct) SetField(name string, value ValueType) {
	if s.fields == nil {
		s.fields = make(map[string]ValueType)
	}
	s.fields[name] = value
}

func (s StructTypeStruct) GetField(name string) (ValueType, bool) {
	value, ok := s.fields[name]
	return value, ok
}

func (s StructTypeStruct) Primitive() ValueType {
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

func (s StructTypeStruct) Equals(other ValueType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(StructTypeStruct); ok {
		if s.Name != other.Name {
			return false
		}
		if len(s.fields) != len(other.fields) {
			return false
		}
		for name, value := range s.fields {
			if otherValue, ok := other.fields[name]; !ok || !value.Equals(otherValue) {
				return false
			}
		}
		return true
	}
	return false
}
