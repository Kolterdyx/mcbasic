package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
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
	return fmt.Sprintf("return fail\n")
}

func (o *Op) GetStructFields(name interfaces.ValueType) []interfaces.StructField {
	for _, s := range o.Structs {
		if s.Name.Lexeme == string(name) {
			return s.Fields
		}
	}
	return nil
}

func Cs(s string) string {
	if s == RCF || s == CALL {
		return s
	}
	if strings.Contains(s, "$(__call__)") {
		return s
	}
	return s + "$(__call__)"
}
