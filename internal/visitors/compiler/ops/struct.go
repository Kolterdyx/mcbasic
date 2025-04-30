package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
)

func (o *Op) StructDefine(structType types.StructTypeStruct) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set value %s\n", o.Namespace, StructPath, structType.Name, o.StructToNbt(structType))
}

func (o *Op) StructToNbt(structType types.StructTypeStruct) string {
	structStmt, ok := o.Structs[structType.Name]
	if !ok {
		log.Errorf("Struct %s not found", structType)
		return ""
	}
	return structStmt.Compound.ToString()
}

func (o *Op) StructGet(from, field, to string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set from storage %s:data %s.%s.%s\n", o.Namespace, VarPath, to, o.Namespace, VarPath, from, field)
}
