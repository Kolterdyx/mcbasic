package types

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	log "github.com/sirupsen/logrus"
)

type PrimitiveType string

const (
	PrimitiveVoidType   PrimitiveType = "void"
	PrimitiveIntType    PrimitiveType = "int"
	PrimitiveStringType PrimitiveType = "str"
	PrimitiveDoubleType PrimitiveType = "double"
)

type PrimitiveTypeStruct struct {
	ValueType

	primitiveType PrimitiveType
}

func (p PrimitiveTypeStruct) Primitive() ValueType {
	return p
}

func (p PrimitiveTypeStruct) ToString() string {
	return string(p.primitiveType)
}

func (p PrimitiveTypeStruct) ToNBT() nbt.Value {
	switch p.primitiveType {
	case PrimitiveIntType:
		return nbt.NewInt(0)
	case PrimitiveStringType:
		return nbt.NewString("")
	case PrimitiveDoubleType:
		return nbt.NewDouble(0)
	case PrimitiveVoidType:
		break
	}
	log.Fatalf("Should be unreachable: %v", p.primitiveType)
	return nil
}

func (p PrimitiveTypeStruct) Equals(other ValueType) bool {
	if other == nil {
		return false
	}
	if other, ok := other.(PrimitiveTypeStruct); ok {
		return p.primitiveType == other.primitiveType
	}
	return false
}

var (
	VoidType   = PrimitiveTypeStruct{primitiveType: PrimitiveVoidType}
	IntType    = PrimitiveTypeStruct{primitiveType: PrimitiveIntType}
	StringType = PrimitiveTypeStruct{primitiveType: PrimitiveStringType}
	DoubleType = PrimitiveTypeStruct{primitiveType: PrimitiveDoubleType}
)
