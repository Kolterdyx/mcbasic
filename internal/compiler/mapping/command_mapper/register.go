package command_mapper

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *CommandMapper) MoveRaw(storageFrom, pathFrom, storageTo, pathTo string) string {
	return fmt.Sprintf("data modify storage %s %s set from storage %s %s\n", storageTo, pathTo, storageFrom, pathFrom)
}

func (c *CommandMapper) Move(from, to string) string {
	return c.MoveRaw(
		fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf("%s.%s", c.VarPath, from),
		fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf("%s.%s", c.VarPath, to),
	)
}

func (c *CommandMapper) MakeConst(value nbt.Value, to string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s set value %s\n", c.Namespace, c.VarPath, to, value.ToString())
}

func (c *CommandMapper) MoveScore(from string, to string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:data %s.%s\n", to, c.Namespace, c.Namespace, c.VarPath, from)
}

func (c *CommandMapper) LoadScore(from string, to string) string {
	return fmt.Sprintf("execute store result storage %s:data %s.%s int 1 run scoreboard players get %s %s\n", c.Namespace, c.VarPath, to, from, c.Namespace)
}

func (c *CommandMapper) IncScore(varName string) string {
	cmd := ""
	cmd += c.MoveScore(varName, varName)
	cmd += fmt.Sprintf("scoreboard players add %s %s 1\n", varName, c.Namespace)
	cmd += c.LoadScore(varName, varName)
	return cmd
}
