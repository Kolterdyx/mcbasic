package mapping

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"strings"
)

type IntermediateMapper struct {
	Mapper
	MapperData
}

func NewIntermediateMapper(namespace string, structs map[string]statements.StructDeclarationStmt) *IntermediateMapper {
	return &IntermediateMapper{
		MapperData: MapperData{
			Namespace: namespace,
			Structs:   structs,
		},
	}
}

func (c *IntermediateMapper) MacroWrapper(argName string) string {
	return "$(" + argName + ")"
}

func (c *IntermediateMapper) MacroLineIndicator(source string) string {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		if strings.Contains(line, "$(") && !(line[0:1] == "$") {
			lines[i] = "$;" + line
		} else {
			lines[i] = " ;" + line
		}
	}
	return strings.Join(lines, "\n")
}

func (c *IntermediateMapper) MakeRegister(regName string) (newRegName string, cmd string) {
	c.registerCount++
	newRegName = fmt.Sprintf("%s%d", regName, c.registerCount)
	cmd = c.MakeConst(nbt.NewInt(0), c.Cs(newRegName))
	return
}

func (c *IntermediateMapper) Return() string {
	return fmt.Sprintf("RETURN\n")
}
