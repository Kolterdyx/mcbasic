package types

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	log "github.com/sirupsen/logrus"
)

type FunctionTypeStruct struct {
	ValueType

	Name string
}

func NewFunctionType(name string) FunctionTypeStruct {
	return FunctionTypeStruct{
		Name: name,
	}
}

func (f FunctionTypeStruct) Primitive() ValueType {
	log.Fatalln("FunctionTypeStruct.Primitive() should not be called")
	return f
}

func (f FunctionTypeStruct) ToString() string {
	return f.Name
}

func (f FunctionTypeStruct) ToNBT() nbt.Value {
	log.Fatalln("FunctionTypeStruct.ToNBT() should not be called")
	return nbt.NewCompound()
}
