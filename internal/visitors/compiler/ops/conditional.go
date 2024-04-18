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
		if line == "" {
			continue
		}
		cmd += fmt.Sprintf("execute %s %s run %s\n", condType, condition, line)
	}
	return cmd
}
