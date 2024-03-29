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
	RA = "ra"
	RB = "rb"

	RX = "rx"

	RCF = "rcf"
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
