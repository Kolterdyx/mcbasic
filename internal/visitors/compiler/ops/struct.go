package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func (o *Op) StructDefine(structType types.StructTypeStruct) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set value %s\n", o.Namespace, StructPath, structType.Name, o.StructToNbt(structType))
}

func (o *Op) StructToNbt(structType types.StructTypeStruct) string {
	nbtData := "{"
	fields := o.GetStructFields(structType)
	if fields == nil {
		log.Errorf("Struct %s not found", structType)
		return ""
	}
	for i, field := range fields {
		nbtData += field.Name + ": "
		switch field.Type {
		case types.IntType:
			nbtData += "0L"
		case types.DoubleType:
			nbtData += "0.0d"
		case types.StringType:
			nbtData += "''"
		default:
			switch reflect.TypeOf(field.Type) {
			case reflect.TypeOf(types.ListTypeStruct{}):
				nbtData += "[]"
			case reflect.TypeOf(types.StructTypeStruct{}):
				nbtData += o.StructToNbt(field.Type.(types.StructTypeStruct))
			}
		}
		if i != len(fields)-1 {
			nbtData += ", "
		}
	}
	nbtData += "}"
	return nbtData
}

func (o *Op) StructGet(from, field, to string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set from storage %s:data %s.%s.%s\n", o.Namespace, VarPath, to, o.Namespace, VarPath, from, field)
}
