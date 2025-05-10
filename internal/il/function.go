package il

import (
	"fmt"
)

func (c *Compiler) Call(funcName string) string {
	ns, fn := splitFunctionName(funcName, c.Namespace)
	return c.inst(Call, fn, ns, fmt.Sprintf("%s.%s", ArgPath, fn))
}
