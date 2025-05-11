package ir

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

const (
	RX = "$RX"
	RA = "$RA"
	RB = "$RB"

	RET  = "$RET"  // Function return value
	RETF = "$RETF" // Early return flag
	CALL = "$CALL"

	VarPath    = "vars"
	ArgPath    = "args"
	StructPath = "structs"

	MaxCallCounter = 65536
)

type Function struct {
	interfaces.Function
	Name string
	Code interfaces.IRCode
}

func NewFunction(name string, code interfaces.IRCode) interfaces.Function {
	nc := NewCode(code.GetNamespace(), code.GetStorage())
	return &Function{
		Name: name,
		Code: nc.Extend(code.(*Code)),
	}
}

func (f *Function) ToString() string {
	return f.Code.ToString()
}

func (f *Function) ToMCFunction() string {
	return f.Code.ToMCCode()
}

func (f *Function) GetName() string {
	return f.Name
}

func (f *Function) GetCode() interfaces.IRCode {
	return f.Code
}

func parseInstruction(parts []string, dpNamespace, storage string) Instruction {
	return Instruction{
		Type:        interfaces.InstructionType(parts[0]),
		Args:        parts[1:],
		DPNamespace: dpNamespace,
		Storage:     storage,
	}
}
