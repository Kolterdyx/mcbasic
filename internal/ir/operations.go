package ir

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	log "github.com/sirupsen/logrus"
)

func (c *Code) XSet(storage, path string, value nbt.Value) interfaces.IRCode {
	return c.addInst(Set, storage, path, value.ToString())
}

func (c *Code) Set(path string, value nbt.Value) interfaces.IRCode {
	return c.XSet(c.storage, path, value)
}

func (c *Code) SetVar(name string, value nbt.Value) interfaces.IRCode {
	return c.Set(c.varPath(name), value)
}

func (c *Code) SetArg(funcName, argName string, value nbt.Value) interfaces.IRCode {
	_, fn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.Set(c.argPath(fn, argName), value)
}

func (c *Code) SetArgs(funcName string, value nbt.Compound) interfaces.IRCode {
	_, fn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.Set(fmt.Sprintf("%s.%s", ArgPath, fn), value)
}

func (c *Code) XCopy(storageFrom, from, storageTo, to string) interfaces.IRCode {
	return c.addInst(Copy, storageFrom, from, storageTo, to)
}

func (c *Code) Copy(from, to string) interfaces.IRCode {
	return c.XCopy(c.storage, from, c.storage, to)
}

func (c *Code) CopyVar(from, to string) interfaces.IRCode {
	return c.Copy(c.varPath(from), c.varPath(to))
}

func (c *Code) CopyArg(varName, funcName, argName string) interfaces.IRCode {
	_, fn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.Copy(c.varPath(varName), c.argPath(fn, argName))
}

func (c *Code) XRemove(storage, path string) interfaces.IRCode {
	return c.addInst(Remove, storage, path)
}

func (c *Code) Remove(path string) interfaces.IRCode {
	return c.XRemove(c.storage, path)
}

func (c *Code) RemoveVar(name string) interfaces.IRCode {
	return c.Remove(c.varPath(name))
}

func (c *Code) RemoveArg(funcName, argName string) interfaces.IRCode {
	_, fn := utils.SplitFunctionName(funcName, c.Namespace)
	return c.Remove(c.argPath(fn, argName))
}

func (c *Code) XLoad(path, score string) interfaces.IRCode {
	return c.addInst(Load, path, score)
}

func (c *Code) Load(path, score string) interfaces.IRCode {
	return c.XLoad(c.varPath(path), score)
}

func (c *Code) Store(score, path string) interfaces.IRCode {
	return c.addInst(Store, score, c.varPath(path))
}

func (c *Code) Score(target string, score *nbt.Int) interfaces.IRCode {
	return c.addInst(Score, target, score.ToString())
}

func (c *Code) MathOp(operator string) interfaces.IRCode {
	return c.addInst(Math, operator)
}

func (c *Code) Ret() interfaces.IRCode {
	return c.addInst(Ret)
}

func (c *Code) Size(source, res string) interfaces.IRCode {
	return c.addInst(Size, c.varPath(source), c.varPath(res))
}

func (c *Code) Func(name string) interfaces.IRCode {
	return c.addInst(Func, name)
}

func (c *Code) AppendSet(listPath string, value nbt.Value) interfaces.IRCode {
	return c.addInst(AppendSet, c.varPath(listPath), value.ToString())
}

func (c *Code) AppendCopy(listPath, valuePath string) interfaces.IRCode {
	return c.addInst(AppendCopy, c.varPath(listPath), c.varPath(valuePath))
}

func (c *Code) MakeIndex(valuePath, res string) interfaces.IRCode {
	c.CopyArg(valuePath, "zzz/path/make_index", "index")
	c.SetArg("zzz/path/make_index", "res", nbt.NewString(c.varPath(res)))
	c.SetArg("zzz/path/make_index", "storage", nbt.NewString(c.storage))
	c.Call("mcb:zzz/path/make_index")
	return c
}

func (c *Code) IntCompare(regRa, regRb string, operator interfaces.TokenType, res string) interfaces.IRCode {
	c.SetVar(res, nbt.NewInt(0))
	c.addInst(Cmp, regRa, utils.CmpOperator(operator), regRb, c.varPath(res))
	return c
}

func (c *Code) DoubleCompare(regRa, regRb string, operator interfaces.TokenType, res string) interfaces.IRCode {
	panic("DoubleCompare has not been implemented yet")
}

func (c *Code) If(condVar string, code interfaces.IRCode) interfaces.IRCode {
	for _, inst := range code.GetInstructions() {
		c.addInst(If, condVar, inst.ToString())
	}
	return c
}

func (c *Code) Unless(condVar string, code interfaces.IRCode) interfaces.IRCode {
	for _, inst := range code.GetInstructions() {
		c.addInst(Unless, condVar, inst.ToString())
	}
	return c
}

func (c *Code) ExceptionString(message string) interfaces.IRCode {
	format := nbt.StringFormat{Color: nbt.Red, Italic: true}
	c.SetArg(fmt.Sprintf("mcb:%s/report", paths.Internal), "text", nbt.NewFormattedString(message, format))
	log.Debugf("format: %s", format.ToString())
	c.Call(fmt.Sprintf("mcb:%s/report", paths.Internal))
	return c
}

func (c *Code) ExceptionFormat(messageParts ...nbt.Value) interfaces.IRCode {
	message := nbt.NewList()
	for _, part := range messageParts {
		message.Add(part)
	}
	c.SetArg(fmt.Sprintf("mcb:%s/report", paths.Internal), "text", message)
	c.Call(fmt.Sprintf("mcb:%s/report", paths.Internal))
	return c
}
