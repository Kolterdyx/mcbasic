package ops

import (
	"fmt"
	"strings"
)

const (
	VarPath = "vars"
	ArgPath = "args"
)

const (
	RA = "$RA"
	RB = "$RB"

	RX  = "$RX"
	RO  = "$RO"
	RET = "$RET"

	RCF  = "$RCF"
	CALL = "$CALL"
)

type Op struct {
	Namespace string
	Scope     string
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

func Cs(s string) string {
	if s == RCF || s == CALL {
		return s
	}
	if strings.Contains(s, "$(__call__)") {
		return s
	}
	return s + "$(__call__)"
}
