package mapping

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
)

const (
	VarPath    = "vars"
	ArgPath    = "args"
	StructPath = "structs"
)

const (
	RA = "$RA"
	RB = "$RB"

	RX  = "$RX"
	RET = "$RET"

	CALL = "$CALL"

	MaxCallCounter = 65536
)

type Mapper interface {
	MacroWrapper(argName string) string
	MacroLineIndicator(source string) string
	Cs(name string) string
	MakeRegister(regName string) (newRegName string, cmd string)

	Trace(storage, path string) string

	Return() string

	NegateNumber(varName string) string
	CompNumbers(cond string, ifcond bool, ra string, rb string, rx string) string
	EqNumbers(ra string, rb string, rx string) string
	NeqNumbers(ra string, rb string, rx string) string
	GtNumbers(ra string, rb string, rx string) string
	GteNumbers(ra string, rb string, rx string) string
	LtNumbers(ra string, rb string, rx string) string
	LteNumbers(ra string, rb string, rx string) string
	EqStrings(ra string, rb string, rx string) string
	NeqStrings(ra string, rb string, rx string) string

	ExecCond(condition string, ifcond bool, source string) string

	Call(funcName string, res string) string
	LoadArgs(funcName string, args map[string]string) string
	LoadArg(funcName, argName string, varName string) string
	LoadArgRaw(funcName, argName string, varName string) string
	LoadArgConst(funcName, argName string, value nbt.Value) string

	AppendList(to, from string) string
	MakeIndex(res, index string) string

	IntAdd(x, y, to string) string
	IntSub(x, y, to string) string
	IntMul(x, y, to string) string
	IntDiv(x, y, to string) string
	IntMod(x, y, to string) string
	DoubleAdd(x, y, to string) string
	DoubleSub(x, y, to string) string
	DoubleMul(x, y, to string) string
	DoubleDiv(x, y, to string) string
	DoubleMod(x, y, to string) string
	DoubleSqrt(x, to string) string
	DoubleSin(x, to string) string
	DoubleCos(x, to string) string
	DoubleTan(x, to string) string
	DoubleAsin(x, to string) string
	DoubleAcos(x, to string) string
	DoubleAtan(x, to string) string
	DoubleRound(x, to string) string
	DoubleFloor(x, to string) string
	DoubleCeil(x, to string) string

	PathGet(obj, path, to string) string
	PathSet(obj, path, valuePath string) string
	PathDelete(obj, path string) string

	MoveRaw(storageFrom, pathFrom, storageTo, pathTo string) string
	Move(from, to string) string
	MakeConst(value nbt.Value, to string) string
	MoveScore(from string, to string) string
	LoadScore(from string, to string) string
	IncScore(varName string) string

	Concat(var1, var2, result string) string
	Size(var1, result string) string
	SliceString(from, start, end, result string) string

	StructDefine(structStmt statements.StructDeclarationStmt) string
	StructGet(from, field, to string) string
	StructSet(from, field, to string) string

	Exception(message string) string
}

type MapperData struct {
	Namespace string
	Scope     string
	Structs   map[string]statements.StructDeclarationStmt

	RX   string
	RA   string
	RB   string
	RET  string
	CALL string

	VarPath    string
	ArgPath    string
	StructPath string

	registerCount int
}

func (m *MapperData) NewRegister(regName string) (newRegName string) {
	m.registerCount++
	newRegName = fmt.Sprintf("$%s%d", regName, m.registerCount)
	return
}
