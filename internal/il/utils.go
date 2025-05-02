package il

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (c *Compiler) inst(instruction InstructionType, args ...string) string {
	return fmt.Sprintf("%s %s\n", instruction, strings.Join(args, " "))
}

func (c *Compiler) error(location interfaces.SourceLocation, message string) {
	log.Errorf("[Position %s] Exception: %s\n", location.ToString(), message)
}

func (c *Compiler) varPath(path string) string {
	if strings.HasPrefix(path, fmt.Sprintf("%s.", VarPath)) {
		return path
	}
	return fmt.Sprintf("%s.%s", VarPath, path)
}

func (c *Compiler) argPath(funcName, arg string) string {
	return fmt.Sprintf("%s.%s.%s", ArgPath, funcName, arg)
}

func (c *Compiler) structPath(path string) string {
	return fmt.Sprintf("%s.%s", StructPath, path)
}

func (c *Compiler) makeReg(reg string) string {
	c.registerCounter++
	return fmt.Sprintf("%s%d", reg, c.registerCounter)
}

// Searches the current scopes for functions and variables, returns the type of the variable or function
func (c *Compiler) getReturnType(name string) types.ValueType {
	for _, identifier := range c.scopes[c.currentScope] {
		if identifier.Name == name {
			return identifier.Type
		}
	}
	return types.VoidType
}

func (c *Compiler) macroLineIdentifier(source string) string {
	lines := strings.Split(source, "\n")
	if len(lines) == 0 {
		return ""
	}
	for i, line := range lines {
		if strings.Contains(line, "$(") && !(line[0:1] == "$") {
			lines[i] = "$" + line
		}
	}
	return strings.Join(lines, "\n")
}
