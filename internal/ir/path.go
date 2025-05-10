package ir

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *Compiler) PathGet(obj, path, to string) (cmd string) {
	cmd += c.SetArgs(
		"internal/path/get",
		nbt.NewCompound().
			Set("res", nbt.NewString(c.varPath(to))).
			Set("storage", nbt.NewString(c.storage)).
			Set("obj", nbt.NewString(c.varPath(obj))),
	)
	cmd += c.CopyArg(path, "internal/path/get", "path")
	cmd += c.Call("mcb:internal/path/get")
	return
}

func (c *Compiler) PathSet(obj, path, valuePath string) (cmd string) {
	cmd += c.SetArgs(
		"internal/path/set",
		nbt.NewCompound().
			Set("value_path", nbt.NewString(c.varPath(valuePath))).
			Set("storage", nbt.NewString(c.storage)).
			Set("obj", nbt.NewString(c.varPath(obj))),
	)
	cmd += c.CopyArg(path, "internal/path/set", "path")
	cmd += c.Call("mcb:internal/path/set")
	return
}

func (c *Compiler) PathDelete(obj, path string) (cmd string) {
	cmd += c.SetArgs(
		"internal/path/delete",
		nbt.NewCompound().
			Set("storage", nbt.NewString(c.storage)).
			Set("obj", nbt.NewString(c.varPath(obj))),
	)
	cmd += c.CopyArg(path, "internal/path/delete", "path")
	cmd += c.Call("mcb:internal/path/delete")
	return
}
