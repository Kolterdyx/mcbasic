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

	RCF  = "rcf"
	CALL = "call"
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

func (o *Op) Exec(command string) string {
	return fmt.Sprintf("%s\n", command)
}

func (o *Op) Inc(reg string) string {
	return fmt.Sprintf("scoreboard players add %s %s 1\n", reg, o.Namespace)
}
