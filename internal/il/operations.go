package il

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

func (c *Compiler) XSet(storage, path string, value nbt.Value) string {
	return c.inst(Set, storage, path, value.ToString())
}

func (c *Compiler) Set(path string, value nbt.Value) string {
	return c.XSet(c.storage, path, value)
}

func (c *Compiler) SetVar(name string, value nbt.Value) string {
	return c.Set(c.varPath(name), value)
}

func (c *Compiler) SetArg(funcName, argName string, value nbt.Value) string {
	return c.Set(c.argPath(funcName, argName), value)
}

func (c *Compiler) XCopy(storageFrom, from, storageTo, to string) string {
	return c.inst(Copy, storageFrom, from, storageTo, to)
}

func (c *Compiler) Copy(from, to string) string {
	return c.XCopy(c.storage, from, c.storage, to)
}

func (c *Compiler) CopyVar(from, to string) (cmd string) {
	return c.Copy(c.varPath(from), c.varPath(to))
}

func (c *Compiler) CopyArg(varName, funcName, argName string) string {
	return c.Copy(c.varPath(varName), c.argPath(funcName, argName))
}

func (c *Compiler) Load(path, score string) string {
	return c.inst(Load, c.varPath(path), score)
}

func (c *Compiler) Store(score, path string) string {
	return c.inst(Store, score, c.varPath(path))
}

func (c *Compiler) Score(target string, score *nbt.Int) string {
	return c.inst(Score, target, score.ToString())
}

func (c *Compiler) MathOp(operator string) string {
	return c.inst(Math, operator)
}

func (c *Compiler) Run(command string) string {
	return c.inst(Run, command)
}

func (c *Compiler) Ret() string {
	return c.inst(Ret)
}

func (c *Compiler) Log(text string) string {
	return c.inst(Log, text)
}

func (c *Compiler) Size(source, res string) string {
	return c.inst(Size, c.varPath(source), c.varPath(res))
}

func (c *Compiler) Func(name string) string {
	return c.inst(Func, name)
}

func (c *Compiler) Append(listPath, valuePath string) string {
	return c.inst(Append, c.varPath(listPath), c.varPath(valuePath))
}

func (c *Compiler) MakeIndex(valuePath, res string) (cmd string) {
	cmd += c.CopyArg(valuePath, "internal/path/make_index", "index")
	cmd += c.CopyArg(res, "internal/path/make_index", "res")
	cmd += c.SetArg("internal/path/make_index", "storage", nbt.NewString(c.storage))
	cmd += c.Call("mcb:internal/path/make_index")
	return
}

func (c *Compiler) IntCompare(regRa, regRb string, operator tokens.TokenType, res string) (cmd string) {
	cmd += c.Load(regRa, regRa)
	cmd += c.Load(regRb, regRb)
	cmd += c.inst(Cmp, regRa, operator.String(), regRb, c.varPath(res))
	return
}

func (c *Compiler) DoubleCompare(regRa, regRb string, operator tokens.TokenType, res string) string {
	panic("implement me")
}

func (c *Compiler) If(condVar, inst string) (cmd string) {
	cmd += c.Load(condVar, condVar)
	cmd += c.inst(If, condVar, inst[:len(inst)-1])
	return
}

func (c *Compiler) Unless(condVar, inst string) (cmd string) {
	cmd += c.inst(Unless, condVar, inst[:len(inst)-1])
	return
}

func (c *Compiler) Exception(message string) (cmd string) {
	cmd += c.SetArg("mcb:error", "text", nbt.NewString(message))
	cmd += c.Call("mcb:error")
	return
}
