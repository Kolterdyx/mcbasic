package ir

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

func (c *Code) StructGet(structPath, field, dataPath string) interfaces.IRCode {
	return c.CopyVar(fmt.Sprintf("%s.%s", structPath, field), dataPath)
}

func (c *Code) StructSet(dataPath, field, structPath string) interfaces.IRCode {
	return c.Copy(c.varPath(dataPath), fmt.Sprintf("%s.%s", c.varPath(structPath), field))
}
