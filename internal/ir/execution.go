package ir

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	"strings"
)

func (c *Code) Call(funcName string) interfaces.IRCode {
	ns, fn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.addInst(Call, fn, ns, fmt.Sprintf("%s.%s", ArgPath, fn))
}

func (c *Code) CallWithArgs(funcName, argPath string) interfaces.IRCode {
	ns, fn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.addInst(Call, fn, ns, argPath)
}

func (c *Code) Branch(branchName, funcName string) interfaces.IRCode {
	bns, bfn := utils.SplitFunctionName(branchName, c.Namespace)
	_, ofn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.addInst(Branch, bfn, bns, fmt.Sprintf("%s.%s", ArgPath, ofn))
}

func (c *Code) Exec(mcCommand string) interfaces.IRCode {
	c.SetArg("mcb:exec", "command", nbt.NewString(strings.TrimSpace(mcCommand)))
	c.SetArg("mcb:exec", "namespace", nbt.NewString(c.Namespace))
	c.Call("mcb:exec")
	return c
}

func (c *Code) ExecReg(varName string) interfaces.IRCode {
	c.CopyArg(varName, "mcb:exec", "command")
	c.SetArg("mcb:exec", "namespace", nbt.NewString(c.Namespace))
	c.Call("mcb:exec")
	return c
}

func (c *Code) Raw(mcCommand string) interfaces.IRCode {
	return c.addInst(Raw, mcCommand)
}
