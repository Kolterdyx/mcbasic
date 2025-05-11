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

func parseInstruction(rawInstruction string, dpNamespace, storage string) Instruction {
	parts := split(rawInstruction)
	return Instruction{
		Type:        interfaces.InstructionType(parts[0]),
		Args:        parts[1:],
		DPNamespace: dpNamespace,
		Storage:     storage,
	}
}

func split(line string) []string {
	// Split the line by spaces, but keep quoted strings together and snbt data together
	parts := make([]string, 0)
	quote := false
	quoteChar := ' '
	brace := 0
	bracket := 0
	start := 0
	for i, r := range line {
		if r == '"' || r == '\'' {
			if quote && r == quoteChar && line[i-1] != '\\' {
				// End of quoted string
				quote = false
			} else if !quote && line[i-1] != '\\' {
				// Start of quoted string
				quote = true
				quoteChar = r
			}
		} else if r == '{' {
			brace++
		} else if r == '}' {
			brace--
		} else if r == '[' {
			bracket++
		} else if r == ']' {
			bracket--
		} else if r == ' ' && !quote && brace == 0 && bracket == 0 {
			if start < i {
				parts = append(parts, line[start:i])
			}
			start = i + 1
		}
	}
	if start < len(line) {
		parts = append(parts, line[start:])
	}
	return parts
}
