package ops

import (
	"fmt"
	"strings"
)

func (o *Op) ExecCond(condition string, ifcond bool, source string) string {
	splitSource := strings.Split(source, "\n")
	cmd := ""
	condType := "if"
	if !ifcond {
		condType = "unless"
	}
	for _, line := range splitSource {
		if line == "" || line[0] == '#' || line[0] == '\n' {
			continue
		}
		cmd += fmt.Sprintf("execute %s %s run %s\n", condType, condition, line)
	}
	return cmd
}
