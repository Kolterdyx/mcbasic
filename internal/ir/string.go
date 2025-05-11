package ir

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (c *Code) StringConcat(a, b, res string) interfaces.IRCode {
	c.SetArgs(
		"internal/string/concat",
		nbt.NewCompound().
			Set("res", nbt.NewString(c.varPath(res))).
			Set("storage", nbt.NewString(c.storage)),
	)
	c.CopyArg(a, "internal/string/concat", "a")
	c.CopyArg(b, "internal/string/concat", "b")
	c.Call("mcb:internal/string/concat")
	return c
}

func (c *Code) StringCompare(a, b, res string) interfaces.IRCode {
	c.SetArgs(
		"internal/string/compare",
		nbt.NewCompound().
			Set("res", nbt.NewString(c.varPath(res))).
			Set("storage", nbt.NewString(c.storage)),
	)
	c.CopyArg(a, "internal/string/compare", "a")
	c.CopyArg(b, "internal/string/compare", "b")
	c.Call("mcb:internal/string/compare")
	return c
}

func (c *Code) StringSlice(stringVar, startIndex, endIndex, res string) interfaces.IRCode {
	c.SetArgs(
		"internal/string/slice",
		nbt.NewCompound().
			Set("storage", nbt.NewString(c.storage)).
			Set("res", nbt.NewString(c.varPath(c.varPath(res)))).
			Set("string", nbt.NewString(c.varPath(stringVar))),
	)
	c.CopyArg(startIndex, "internal/string/slice", "start")
	c.CopyArg(endIndex, "internal/string/slice", "end")
	c.Call("mcb:internal/string/slice")
	return c
}
