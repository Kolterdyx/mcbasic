package command_mapper

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *CommandMapper) Concat(var1, var2, result string) string {
	cmd := ""
	cmd += c.LoadArgConst("internal/concat", "res", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, c.RET)))
	cmd += c.LoadArgConst("internal/concat", "storage", nbt.NewString(fmt.Sprintf("%s:data", c.Namespace)))
	cmd += c.LoadArg("internal/concat", "a", var1)
	cmd += c.LoadArg("internal/concat", "b", var2)
	cmd += c.Call("mcb:internal/concat", result)
	return cmd
}

func (c *CommandMapper) Size(var1, result string) string {
	return fmt.Sprintf("execute store result storage %s:data %s.%s int 1 run data get storage %s:data %s.%s\n", c.Namespace, c.VarPath, result, c.Namespace, c.VarPath, var1)
}

func (c *CommandMapper) SliceString(from, start, end, result string) string {
	cmd := ""
	cmd += c.LoadArgConst("internal/slice", "res", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, c.RET)))
	cmd += c.LoadArgConst("internal/slice", "storage", nbt.NewString(fmt.Sprintf("%s:data", c.Namespace)))
	cmd += c.LoadArgConst("internal/slice", "from", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, from)))
	cmd += c.LoadArg("internal/slice", "start", start)
	cmd += c.LoadArg("internal/slice", "end", end)
	cmd += c.Call("mcb:internal/slice", result)
	return cmd
}
