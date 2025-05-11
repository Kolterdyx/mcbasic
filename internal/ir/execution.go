package ir

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/utils"
)

func (c *Code) Call(funcName string) interfaces.IRCode {
	ns, fn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.addInst(Call, fn, ns, fmt.Sprintf("%s.%s", ArgPath, fn))
}

func (c *Code) Branch(branchName, funcName string) interfaces.IRCode {
	bns, bfn := utils.SplitFunctionName(branchName, c.Namespace)
	_, ofn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.addInst(Branch, bfn, bns, fmt.Sprintf("%s.%s", ArgPath, ofn))
}

func (c *Code) IBranch(branchName, funcName string) interfaces.Instruction {
	bns, bfn := utils.SplitFunctionName(branchName, c.Namespace)
	_, ofn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.makeInst(Branch, bfn, bns, fmt.Sprintf("%s.%s", ArgPath, ofn))
}

func (c *Code) Exec(mcCommand string) interfaces.IRCode {
	c.SetArg("mcb:exec", "command", nbt.NewString(mcCommand))
	c.Call("mcb:exec")
	return c
}
