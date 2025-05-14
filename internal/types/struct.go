package types

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/elliotchance/orderedmap/v3"
	"slices"
)

type StructTypeStruct struct {
	ValueType

	Name   string
	fields *orderedmap.OrderedMap[string, ValueType]
}

func NewStructType(name string) StructTypeStruct {
	return StructTypeStruct{
		Name:   name,
		fields: orderedmap.NewOrderedMap[string, ValueType](),
	}
}

func (s StructTypeStruct) SetField(name string, value ValueType) {
	s.fields.Set(name, value)
}

func (s StructTypeStruct) GetFieldType(name string) (ValueType, bool) {
	return s.fields.Get(name)
}

func (s StructTypeStruct) GetFieldNames() []string {
	return slices.Collect(s.fields.Keys())
}

func (s StructTypeStruct) Primitive() ValueType {
	return s
}

func (s StructTypeStruct) ToString() string {
	return s.Name
}

func (s StructTypeStruct) ToNBT() nbt.Value {
	compound := nbt.NewCompound()
	for name, field := range s.fields.AllFromFront() {
		compound.Set(name, field.ToNBT())
	}
	return compound
}

func (s StructTypeStruct) Size() int {
	return s.fields.Len()
}

func (s StructTypeStruct) Equals(other ValueType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(StructTypeStruct); ok {
		if s.Name != other.Name {
			return false
		}
		if s.fields.Len() != other.fields.Len() {
			return false
		}
		for name, field := range s.fields.AllFromFront() {
			if otherField, ok := other.GetFieldType(name); !ok || !field.Equals(otherField) {
				return false
			}
		}
		return true
	}
	return false
}
