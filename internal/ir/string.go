package ir

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *Code) StringConcat(a, b, res string) interfaces.IRCode {
	c.SetArgs(
		"zzz/string/concat",
		nbt.NewCompound().
			Set("res", nbt.NewString(c.varPath(res))).
			Set("storage", nbt.NewString(c.storage)),
	)
	c.CopyArg(a, "zzz/string/concat", "a")
	c.CopyArg(b, "zzz/string/concat", "b")
	c.Call("mcb:zzz/string/concat")
	return c
}

func (c *Code) StringCompare(a, b, res string) interfaces.IRCode {
	c.SetArgs(
		"zzz/string/compare",
		nbt.NewCompound().
			Set("res", nbt.NewString(c.varPath(res))).
			Set("storage", nbt.NewString(c.storage)),
	)
	c.CopyArg(a, "zzz/string/compare", "a")
	c.CopyArg(b, "zzz/string/compare", "b")
	c.Call("mcb:zzz/string/compare")
	return c
}

func (c *Code) StringSlice(stringVar, startIndex, endIndex, res string) interfaces.IRCode {
	c.SetArgs(
		"zzz/string/slice",
		nbt.NewCompound().
			Set("storage", nbt.NewString(c.storage)).
			Set("res", nbt.NewString(c.varPath(c.varPath(res)))).
			Set("string", nbt.NewString(c.varPath(stringVar))),
	)
	c.CopyArg(startIndex, "zzz/string/slice", "start")
	c.CopyArg(endIndex, "zzz/string/slice", "end")
	c.Call("mcb:zzz/string/slice")
	return c
}
