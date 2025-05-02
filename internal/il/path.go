package il

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *Compiler) PathGet(obj, path, to string) (cmd string) {
	cmd += c.Set(c.argPath("internal/path/get", "res"), nbt.NewString(c.varPath(to)))
	cmd += c.Set(c.argPath("internal/path/get", "storage"), nbt.NewString(c.storage))
	cmd += c.Set(c.argPath("internal/path/get", "obj"), nbt.NewString(c.varPath(obj)))
	cmd += c.Copy(path, c.argPath("internal/path/get", "path"))
	cmd += c.Call("mcb:internal/path/get")
	return
}

func (c *Compiler) PathSet(obj, path, valuePath string) (cmd string) {
	cmd += c.Set(c.argPath("internal/path/set", "value_path"), nbt.NewString(c.varPath(valuePath)))
	cmd += c.Set(c.argPath("internal/path/set", "storage"), nbt.NewString(c.storage))
	cmd += c.Set(c.argPath("internal/path/set", "obj"), nbt.NewString(c.varPath(obj)))
	cmd += c.Copy(path, c.argPath("internal/path/set", "path"))
	cmd += c.Call("mcb:internal/path/set")
	return
}

func (c *Compiler) PathDelete(obj, path string) (cmd string) {
	cmd += c.Set(c.argPath("internal/path/delete", "storage"), nbt.NewString(c.storage))
	cmd += c.Set(c.argPath("internal/path/delete", "obj"), nbt.NewString(c.varPath(obj)))
	cmd += c.Copy(path, c.argPath("internal/path/delete", "path"))
	cmd += c.Call("mcb:internal/path/delete")
	return
}
