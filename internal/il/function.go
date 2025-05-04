package il

import (
	"fmt"
)

func (c *Compiler) Call(funcName, ret string) string {
	ns, fn := splitFunctionName(funcName, c.Namespace)
	if ret == "" {
		ret = RET
	}
	return c.inst(Call, fn, ns, fmt.Sprintf("%s.%s", ArgPath, fn), ret)
}
