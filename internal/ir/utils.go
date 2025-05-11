package ir

import (
	"fmt"
	"strings"
)

func (c *Code) varPath(path string) string {
	if strings.HasPrefix(path, fmt.Sprintf("%s.", VarPath)) {
		return path
	}
	return fmt.Sprintf("%s.%s.$(__call__)", VarPath, path)
}

func (c *Code) argPath(funcName, arg string) string {
	return fmt.Sprintf("%s.%s.%s", ArgPath, funcName, arg)
}

func (c *Code) structPath(structName string) string {
	return fmt.Sprintf("%s.%s", StructPath, structName)
}
