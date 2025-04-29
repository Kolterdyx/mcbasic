package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"strings"
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

	RCF  = "$RCF"
	CALL = "$CALL"
)

type Op struct {
	Namespace string
	Scope     string
	Structs   []statements.StructDeclarationStmt
}

func (o *Op) Macro(argName string) string {
	return fmt.Sprintf("$(%s)", argName)
}

// MacroReplace add $ at the start of each line that uses macros. Macros are found in the pattern $(name)
func (o *Op) MacroReplace(source string) string {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		if strings.Contains(line, "$(") {
			lines[i] = "$" + line
		}
	}
	return strings.Join(lines, "\n")
}

func (o *Op) Return() string {
	return fmt.Sprintf("return 1\n")
}

func (o *Op) GetStructFields(structType types.StructTypeStruct) []interfaces.StructField {
	for _, s := range o.Structs {
		if s.Name.Lexeme == structType.Name {
			return s.Fields
		}
	}
	return nil
}

func Cs(s string) string {
	if s == RCF || s == CALL {
		return s
	}
	suffix := ".$(__call__)"
	if strings.Contains(s, suffix) {
		return s
	}
	return s + suffix
}

func (o *Op) Trace(name string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/trace", "storage", fmt.Sprintf("%s:data", o.Namespace))
	cmd += o.LoadArgConst("internal/trace", "path", fmt.Sprintf("%s:%s", VarPath, name))
	cmd += o.LoadArg("internal/trace", "value", name)
	cmd += o.Call("mcb:internal/trace", "")
	return cmd
}

func (o *Op) TraceRaw(name string) string {
	cmd := ""
	cmd += o.LoadArgRaw("internal/trace", "value", name)
	cmd += o.Call("mcb:internal/trace", "")
	return cmd
}
