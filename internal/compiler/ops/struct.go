package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
)

func (o *Op) StructDefine(structStmt statements.StructDeclarationStmt) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set value %s\n", o.Namespace, StructPath, structStmt.Name.Lexeme, structStmt.StructType.ToNBT().ToString())
}

func (o *Op) StructGet(from, field, to string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set from storage %s:data %s.%s.%s\n", o.Namespace, VarPath, to, o.Namespace, VarPath, from, field)
}

func (o *Op) StructSet(from, field, to string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set from storage %s:data %s.%s\n", o.Namespace, VarPath, to, field, o.Namespace, VarPath, from)
}
