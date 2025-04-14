package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	log "github.com/sirupsen/logrus"
)

func (o *Op) StructDefine(name string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set value %s\n", o.Namespace, StructPath, name, o.StructToNbt(interfaces.ValueType(name)))
}

func (o *Op) StructToNbt(structName interfaces.ValueType) string {
	nbtData := "{"
	fields := o.GetStructFields(structName)
	if fields == nil {
		log.Errorf("Struct %s not found", structName)
		return ""
	}
	for i, field := range fields {
		nbtData += field.Name + ": "
		switch field.Type {
		case expressions.IntType:
			nbtData += "0L"
		case expressions.DoubleType:
			nbtData += "0.0d"
		case expressions.StringType:
			nbtData += "''"
		case expressions.ListIntType:
			fallthrough
		case expressions.ListDoubleType:
			fallthrough
		case expressions.ListStringType:
			nbtData += "[]"
		default:
			nbtData += o.StructToNbt(field.Type)
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
