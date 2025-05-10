package ir

import (
	"fmt"
)

func (c *Compiler) Call(funcName string) string {
	ns, fn := splitFunctionName(funcName, c.Namespace)
	return c.inst(Call, fn, ns, fmt.Sprintf("%s.%s", ArgPath, fn))
}

func (c *Compiler) Branch(branchName, funcName string) string {
	bns, bfn := splitFunctionName(branchName, c.Namespace)
	_, ofn := splitFunctionName(funcName, c.Namespace)
	return c.inst(Branch, bfn, bns, fmt.Sprintf("%s.%s", ArgPath, ofn))
}
