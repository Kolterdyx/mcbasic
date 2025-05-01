package command_mapper

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *CommandMapper) AppendList(to, from string) string {
	return fmt.Sprintf("data modify storage example:data %s.%s append from storage %s:data %s.%s\n", c.VarPath, to, c.Namespace, c.VarPath, from)
}

func (c *CommandMapper) MakeIndex(res, index string) string {
	cmd := ""
	cmd += c.LoadArgConst("internal/path/make_index", "storage", nbt.NewString(fmt.Sprintf("%s:data", c.Namespace)))
	cmd += c.LoadArgConst("internal/path/make_index", "res", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, c.RET)))
	cmd += c.LoadArg("internal/path/make_index", "index", index)
	cmd += c.Call("mcb:internal/path/make_index", res)
	return cmd
}
