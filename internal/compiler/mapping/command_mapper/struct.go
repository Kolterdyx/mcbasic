package command_mapper

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
)

func (c *CommandMapper) StructDefine(structStmt statements.StructDeclarationStmt) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set value %s\n", c.Namespace, c.StructPath, structStmt.Name.Lexeme, structStmt.StructType.ToNBT().ToString())
}

func (c *CommandMapper) StructGet(from, field, to string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set from storage %s:data %s.%s.%s\n", c.Namespace, c.VarPath, to, c.Namespace, c.VarPath, from, field)
}

func (c *CommandMapper) StructSet(from, field, to string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set from storage %s:data %s.%s\n", c.Namespace, c.VarPath, to, field, c.Namespace, c.VarPath, from)
}
