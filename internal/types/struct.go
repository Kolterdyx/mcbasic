package types

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type StructTypeStruct struct {
	interfaces.ValueType

	Name   string
	Fields map[string]interfaces.ValueType
}

func NewStructType(name string) StructTypeStruct {
	return StructTypeStruct{
		Name:   name,
		Fields: make(map[string]interfaces.ValueType),
	}
}

func (s StructTypeStruct) SetField(name string, value interfaces.ValueType) {
	if s.Fields == nil {
		s.Fields = make(map[string]interfaces.ValueType)
	}
	s.Fields[name] = value
}

func (s StructTypeStruct) GetField(name string) (interfaces.ValueType, bool) {
	value, ok := s.Fields[name]
	return value, ok
}

func (s StructTypeStruct) Primitive() interfaces.ValueType {
	return s
}

func (s StructTypeStruct) ToString() string {
	return s.Name
}
