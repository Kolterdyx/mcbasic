package ir

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/utils"
)

const (
	TraceStorage = "storage"
	TraceScore   = "score"
)

func (c *Code) XTrace(name, storage, path string) interfaces.IRCode {
	return c.addInst(Trace, TraceStorage, name, storage, path)
}

func (c *Code) Trace(name, path string) interfaces.IRCode {
	return c.XTrace(name, c.storage, path)
}

func (c *Code) TraceVar(name, varName string) interfaces.IRCode {
	return c.Trace(name, c.varPath(varName))
}

func (c *Code) TraceArg(name, funcName, argName string) interfaces.IRCode {
	_, fn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.Trace(name, c.argPath(fn, argName))
}

func (c *Code) TraceScore(name, score string) interfaces.IRCode {
	return c.addInst(Trace, TraceScore, name, score)
}
