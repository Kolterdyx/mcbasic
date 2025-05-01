package command_mapper

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/compiler/mapping"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"strings"
)

type CommandMapper struct {
	mapping.Mapper
	mapping.MapperData
}

func NewCommandMapper(namespace string, structs map[string]statements.StructDeclarationStmt) *CommandMapper {
	return &CommandMapper{
		MapperData: mapping.MapperData{
			Namespace:  namespace,
			Structs:    structs,
			RX:         mapping.RX,
			RA:         mapping.RA,
			RB:         mapping.RB,
			RET:        mapping.RET,
			CALL:       mapping.CALL,
			VarPath:    mapping.VarPath,
			ArgPath:    mapping.ArgPath,
			StructPath: mapping.StructPath,
		},
	}
}

func (c *CommandMapper) Cs(s string) string {
	if s == mapping.CALL {
		return s
	}
	suffix := ".$(__call__)"
	if strings.Contains(s, suffix) {
		return s
	}
	return s + suffix
}

func (c *CommandMapper) MacroWrapper(argName string) string {
	return fmt.Sprintf("$(%s)", argName)
}

func (c *CommandMapper) MacroLineIndicator(source string) string {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		if strings.Contains(line, "$(") && !(line[0:1] == "$") {
			lines[i] = "$" + line
		}
	}
	return strings.Join(lines, "\n")
}

func (c *CommandMapper) MakeRegister(regName string) (newRegName string, cmd string) {
	newRegName = c.NewRegister(regName)
	cmd = c.MakeConst(nbt.NewInt(0), c.Cs(newRegName))
	return
}

func (c *CommandMapper) Return() string {
	return fmt.Sprintf("return 1\n")
}

func (c *CommandMapper) Trace(storage, path string) string {

	text := nbt.NewList(
		nbt.NewCompound().
			Set("text", nbt.NewString(path+": ")).
			Set("color", nbt.NewString("green")),
		nbt.NewCompound().
			Set("type", nbt.NewAny("nbt")).
			Set("source", nbt.NewAny("storage")).
			Set("nbt", nbt.NewString(path)).
			Set("storage", nbt.NewString(storage)).
			Set("color", nbt.NewString("yellow")),
	)

	return fmt.Sprintf("tellraw @a[tag=mcblog] %s\n", text.ToString())
}
