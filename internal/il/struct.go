package il

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
)

func (c *Compiler) StructDefine(structStmt statements.StructDeclarationStmt) string {
	return c.Set(c.structPath(structStmt.Name.Lexeme), structStmt.StructType.ToNBT())
}

func (c *Compiler) StructGet(structPath, field, dataPath string) string {
	return c.CopyVar(fmt.Sprintf("%s.%s", structPath, field), dataPath)
}

func (c *Compiler) StructSet(dataPath, field, structPath string) string {
	return c.CopyVar(dataPath, fmt.Sprintf("%s.%s", structPath, field))
}
