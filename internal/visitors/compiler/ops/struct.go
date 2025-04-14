package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

func (o *Op) StructDefine(name string, fields []interfaces.StructField) string {
	nbtData := "{"
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
			// nothing
		}
		if i != len(fields)-1 {
			nbtData += ", "
		}
	}
	nbtData += "}"
	return fmt.Sprintf("data modify storage %s:data %s.%s set value %s\n", o.Namespace, StructPath, name, nbtData)
}
