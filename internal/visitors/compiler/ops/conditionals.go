package ops

import (
	"fmt"
	"strings"
)

// Ift only runs the command if the condition is true.
func (o *Op) Ift(condition string, commands []string) string {
	cmd := ""
	for _, c := range commands {
		if c == "" {
			continue
		}
		cmd += fmt.Sprintf("execute if %s run %s\n", condition, c)
	}
	return cmd
}

// Ifn only runs the command if the condition is false.
func (o *Op) Ifn(condition string, commands []string) string {
	cmd := ""
	for _, c := range commands {
		if c == "" {
			continue
		}
		cmd += fmt.Sprintf("execute unless %s run %s\n", condition, c)
	}
	return cmd
}

// Iftn runs the first command if the condition is true, otherwise it runs the second command.
func (o *Op) Iftn(condition string, trueBranch []string, falseBranch []string) string {
	return o.Ift(condition, trueBranch) + o.Ifn(condition, falseBranch)
}

// Ifc is an alias for an indefinitely long chain of if-elseif-else statements.
func (o *Op) Ifc(statements map[string][]string) string {
	cmd := ""
	i := 0
	prev := ""
	for cond, s := range statements {
		// each subsequent condition must check if the previous condition was false
		tmp := o.Ift(cond, s)
		if i > 0 {
			tmp = o.Ifn(prev, strings.Split(tmp, "\n"))
		}
		prev = cond
		cmd += tmp
		i++
	}
	return cmd
}
