package command_mapper

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *CommandMapper) PathGet(obj, path, to string) string {
	cmd := "### BEGIN Path get operation  ###\n"
	cmd += c.LoadArgConst("internal/path/get", "res", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, c.RET)))
	cmd += c.LoadArgConst("internal/path/get", "storage", nbt.NewString(fmt.Sprintf("%s:data", c.Namespace)))
	cmd += c.LoadArgConst("internal/path/get", "obj", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, obj)))
	cmd += c.LoadArg("internal/path/get", "path", path)
	cmd += c.Call("mcb:internal/path/get", to)
	cmd += "### END   Path get operation ###\n"
	return cmd
}

func (c *CommandMapper) PathSet(obj, path, valuePath string) string {
	cmd := ""
	cmd += c.LoadArgConst("internal/path/set", "storage", nbt.NewString(fmt.Sprintf("%s:data", c.Namespace)))
	cmd += c.LoadArgConst("internal/path/set", "obj", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, obj)))
	cmd += c.LoadArgConst("internal/path/set", "value_path", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, valuePath)))
	cmd += c.LoadArg("internal/path/set", "path", path)
	cmd += c.Call("mcb:internal/path/set", "")
	return cmd
}

func (c *CommandMapper) PathDelete(obj, path string) string {
	cmd := ""
	cmd += c.LoadArgConst("internal/path/delete", "storage", nbt.NewString(fmt.Sprintf("%s:data", c.Namespace)))
	cmd += c.LoadArgConst("internal/path/delete", "obj", nbt.NewString(fmt.Sprintf("%s.%s", c.VarPath, obj)))
	cmd += c.LoadArg("internal/path/delete", "path", path)
	cmd += c.Call("mcb:internal/path/delete", "")
	return cmd
}
