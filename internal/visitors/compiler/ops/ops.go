package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
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
	Structs   map[string]statements.StructDeclarationStmt
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

func (o *Op) Trace(path string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/trace", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/trace", "path", nbt.NewString(path))
	cmd += o.LoadArg("internal/trace", "value", path)
	cmd += o.Call("mcb:internal/trace", "")
	return cmd
}

func (o *Op) TraceRaw(path string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/trace", "storage", nbt.NewString("N/A"))
	cmd += o.LoadArgConst("internal/trace", "path", nbt.NewString(path))
	cmd += o.LoadArgRaw("internal/trace", "value", path)
	cmd += o.Call("mcb:internal/trace", "")
	return cmd
}
