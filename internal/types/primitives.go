package types

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	log "github.com/sirupsen/logrus"
)

const (
	// PrimitiveErrorType  used for error handling
	PrimitiveErrorType interfaces.PrimitiveType = "error"

	PrimitiveVoidType   interfaces.PrimitiveType = "void"
	PrimitiveIntType    interfaces.PrimitiveType = "int"
	PrimitiveStringType interfaces.PrimitiveType = "str"
	PrimitiveDoubleType interfaces.PrimitiveType = "double"
)

type PrimitiveTypeStruct struct {
	interfaces.ValueType

	primitiveType interfaces.PrimitiveType
}

func (p PrimitiveTypeStruct) Primitive() interfaces.ValueType {
	switch p.primitiveType {
	case PrimitiveVoidType:
		return VoidType
	case PrimitiveIntType:
		return IntType
	case PrimitiveStringType:
		return StringType
	case PrimitiveDoubleType:
		return DoubleType
	case PrimitiveErrorType:
		return ErrorType
	}
	log.Fatalf("Should be unreachable: %v", p.primitiveType)
	return nil
}

func (p PrimitiveTypeStruct) ToString() string {
	switch p.primitiveType {
	case PrimitiveVoidType:
		return "void"
	case PrimitiveIntType:
		return "int"
	case PrimitiveStringType:
		return "str"
	case PrimitiveDoubleType:
		return "double"
	case PrimitiveErrorType:
		return "error"
	}
	log.Fatalf("Should be unreachable: %v", p.primitiveType)
	return ""
}

func (p PrimitiveTypeStruct) ToNBT() nbt.Value {
	switch p.primitiveType {
	case PrimitiveIntType:
		return nbt.NewInt(0)
	case PrimitiveStringType:
		return nbt.NewString("")
	case PrimitiveDoubleType:
		return nbt.NewDouble(0)
	case PrimitiveErrorType:
		fallthrough
	case PrimitiveVoidType:
		break
	}
	log.Fatalf("Should be unreachable: %v", p.primitiveType)
	return nil
}

var (
	VoidType   = PrimitiveTypeStruct{primitiveType: PrimitiveVoidType}
	IntType    = PrimitiveTypeStruct{primitiveType: PrimitiveIntType}
	StringType = PrimitiveTypeStruct{primitiveType: PrimitiveStringType}
	DoubleType = PrimitiveTypeStruct{primitiveType: PrimitiveDoubleType}
	ErrorType  = PrimitiveTypeStruct{primitiveType: PrimitiveErrorType}
)
