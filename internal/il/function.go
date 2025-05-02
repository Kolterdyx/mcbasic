package il

import (
	"fmt"
	"strings"
)

func (c *Compiler) Call(funcName string) string {
	if strings.Contains(funcName, ":") {
		return c.inst(Call, funcName)
	} else {
		return c.inst(Call, fmt.Sprintf("%s:%s", c.Namespace, funcName))
	}
}
