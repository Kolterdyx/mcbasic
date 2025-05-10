package ir

import "github.com/Kolterdyx/mcbasic/internal/nbt"

func (c *Compiler) StringConcat(a, b, res string) (cmd string) {
	cmd += c.SetArgs(
		"internal/string/concat",
		nbt.NewCompound().
			Set("res", nbt.NewString(res)).
			Set("storage", nbt.NewString(c.storage)),
	)
	cmd += c.CopyArg(a, "internal/string/concat", "a")
	cmd += c.CopyArg(b, "internal/string/concat", "b")
	cmd += c.Call("mcb:internal/string/concat")
	return
}

func (c *Compiler) StringCompare(a, b, res string) (cmd string) {
	cmd += c.SetArgs(
		"internal/string/compare",
		nbt.NewCompound().
			Set("res", nbt.NewString(res)).
			Set("storage", nbt.NewString(c.storage)),
	)
	cmd += c.CopyArg(a, "internal/string/compare", "a")
	cmd += c.CopyArg(b, "internal/string/compare", "b")
	cmd += c.Call("mcb:internal/string/compare")
	return
}

func (c *Compiler) StringSlice(stringVar, startIndex, endIndex, res string) (cmd string) {
	cmd += c.SetArgs(
		"internal/string/slice",
		nbt.NewCompound().
			Set("storage", nbt.NewString(c.storage)).
			Set("res", nbt.NewString(c.varPath(res))).
			Set("string", nbt.NewString(c.varPath(stringVar))),
	)
	cmd += c.CopyArg(startIndex, "internal/string/slice", "start")
	cmd += c.CopyArg(endIndex, "internal/string/slice", "end")
	cmd += c.Call("mcb:internal/string/slice")
	return
}
