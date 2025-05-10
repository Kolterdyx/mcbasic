package ir

import "github.com/Kolterdyx/mcbasic/internal/nbt"

func (c *Compiler) StringConcat(a, b, res string) (cmd string) {
	cmd += c.CopyArg(a, "internal/string/concat", "a")
	cmd += c.CopyArg(b, "internal/string/concat", "b")
	cmd += c.SetArg("internal/string/concat", "res", nbt.NewString(c.varPath(res)))
	cmd += c.SetArg("internal/string/concat", "storage", nbt.NewString(c.storage))
	cmd += c.Call("mcb:internal/string/concat")
	return
}

func (c *Compiler) StringCompare(a, b, res string) (cmd string) {
	cmd += c.CopyArg(a, "internal/string/compare", "a")
	cmd += c.CopyArg(b, "internal/string/compare", "b")
	cmd += c.SetArg("internal/string/compare", "res", nbt.NewString(res))
	cmd += c.SetArg("internal/string/compare", "storage", nbt.NewString(c.storage))
	cmd += c.Call("mcb:internal/string/compare")
	return
}

func (c *Compiler) StringSlice(stringVar, startIndex, endIndex, res string) (cmd string) {
	cmd += c.SetArg("internal/string/slice", "storage", nbt.NewString(c.storage))
	cmd += c.SetArg("internal/string/slice", "res", nbt.NewString(c.varPath(res)))
	cmd += c.SetArg("internal/string/slice", "string", nbt.NewString(c.varPath(stringVar)))
	cmd += c.CopyArg(startIndex, "internal/string/slice", "start")
	cmd += c.CopyArg(endIndex, "internal/string/slice", "end")
	cmd += c.Call("mcb:internal/string/slice")
	return
}
