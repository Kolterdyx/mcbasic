package ir

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *Code) PathGet(obj, path, to string) interfaces.IRCode {
	c.SetArgs(
		"zzz/path/get",
		nbt.NewCompound().
			Set("res", nbt.NewString(c.varPath(to))).
			Set("storage", nbt.NewString(c.storage)).
			Set("obj", nbt.NewString(c.varPath(obj))),
	)
	c.CopyArg(path, "zzz/path/get", "path")
	c.Call("mcb:zzz/path/get")
	return c
}

func (c *Code) PathSet(obj, path, valuePath string) interfaces.IRCode {
	c.SetArgs(
		"zzz/path/set",
		nbt.NewCompound().
			Set("value_path", nbt.NewString(c.varPath(valuePath))).
			Set("storage", nbt.NewString(c.storage)).
			Set("obj", nbt.NewString(c.varPath(obj))),
	)
	c.CopyArg(path, "zzz/path/set", "path")
	c.Call("mcb:zzz/path/set")
	return c
}

func (c *Code) PathDelete(obj, path string) interfaces.IRCode {
	c.SetArgs(
		"zzz/path/delete",
		nbt.NewCompound().
			Set("storage", nbt.NewString(c.storage)).
			Set("obj", nbt.NewString(c.varPath(obj))),
	)
	c.CopyArg(path, "zzz/path/delete", "path")
	c.Call("mcb:zzz/path/delete")
	return c
}
