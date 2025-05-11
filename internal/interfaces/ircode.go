package interfaces

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

type InstructionType string

type Instruction interface {
	ToString() string
	ToMCCommand() string

	GetType() InstructionType
	GetArgs() []string
}

type IRCode interface {
	ToString() string
	ToMCCode() string
	Len() int
	Extend(IRCode) IRCode
	AddInstruction(Instruction) IRCode

	GetInstructions() []Instruction
	SetInstructions([]Instruction)
	GetNamespace() string
	GetStorage() string

	Call(funcName string) IRCode
	Branch(branchName, funcName string) IRCode
	Exec(mcCommand string) IRCode
	Raw(mcCommand string) IRCode

	IntAdd(x, y, to string) IRCode
	IntSub(x, y, to string) IRCode
	IntMul(x, y, to string) IRCode
	IntDiv(x, y, to string) IRCode
	IntMod(x, y, to string) IRCode
	DoubleAdd(x, y, to string) IRCode
	DoubleSub(x, y, to string) IRCode
	DoubleMul(x, y, to string) IRCode
	DoubleDiv(x, y, to string) IRCode
	DoubleMod(x, y, to string) IRCode
	DoubleSqrt(x, to string) IRCode
	DoubleCos(x, to string) IRCode
	DoubleSin(x, to string) IRCode
	DoubleTan(x, to string) IRCode
	DoubleAcos(x, to string) IRCode
	DoubleAsin(x, to string) IRCode
	DoubleAtan(x, to string) IRCode
	DoubleAtan2(x, y, to string) IRCode
	DoubleFloor(x, to string) IRCode
	DoubleCeil(x, to string) IRCode
	DoubleRound(x, to string) IRCode

	XSet(storage, path string, value nbt.Value) IRCode
	Set(path string, value nbt.Value) IRCode
	SetVar(name string, value nbt.Value) IRCode
	SetArg(funcName, argName string, value nbt.Value) IRCode

	SetArgs(funcName string, value nbt.Compound) IRCode
	XCopy(storageFrom, from, storageTo, to string) IRCode
	Copy(from, to string) IRCode
	CopyVar(from, to string) IRCode
	CopyArg(varName, funcName, argName string) IRCode
	XRemove(storage, path string) IRCode
	Remove(path string) IRCode
	RemoveVar(name string) IRCode
	RemoveArg(funcName, argName string) IRCode
	XLoad(path, score string) IRCode
	Load(path, score string) IRCode
	Store(score, path string) IRCode
	Score(target string, score *nbt.Int) IRCode
	MathOp(operator string) IRCode
	Ret() IRCode
	Size(source, res string) IRCode
	Func(name string) IRCode
	AppendSet(listPath string, value nbt.Value) IRCode
	AppendCopy(listPath, valuePath string) IRCode
	MakeIndex(valuePath, res string) IRCode
	IntCompare(regRa, regRb string, operator TokenType, res string) IRCode
	DoubleCompare(regRa, regRb string, operator TokenType, res string) IRCode
	If(condVar string, code IRCode) IRCode
	Unless(condVar string, code IRCode) IRCode
	Exception(message string) IRCode

	PathGet(obj, path, to string) IRCode
	PathSet(obj, path, valuePath string) IRCode
	PathDelete(obj, path string) IRCode

	StringConcat(a, b, res string) IRCode
	StringCompare(a, b, res string) IRCode
	StringSlice(stringVar, startIndex, endIndex, res string) IRCode

	StructGet(structPath, field, dataPath string) IRCode
	StructSet(dataPath, field, structPath string) IRCode

	XTrace(name, storage, path string) IRCode
	Trace(name, path string) IRCode
	TraceVar(name, varName string) IRCode
	TraceArg(name, funcName, argName string) IRCode
	TraceScore(name, score string) IRCode
}

type Function interface {
	GetName() string
	GetCode() IRCode
	ToString() string
	ToMCFunction() string
}
