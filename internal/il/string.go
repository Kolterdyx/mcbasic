package il

import "github.com/Kolterdyx/mcbasic/internal/nbt"

func (c *Compiler) StringConcat(a, b, res string) (cmd string) {
	cmd += c.Copy(c.varPath(a), c.argPath("internal/string/concat", "a"))
	cmd += c.Copy(c.varPath(b), c.argPath("internal/string/concat", "b"))
	cmd += c.SetArg("internal/string/concat", "res", nbt.NewString(res))
	cmd += c.SetArg("internal/string/concat", "storage", nbt.NewString(c.storage))
	cmd += c.Call("mcb:internal/string/concat")
	return
}

func (c *Compiler) StringCompare(a, b, res string) (cmd string) {
	cmd += c.Copy(c.varPath(a), c.argPath("internal/string/compare", "a"))
	cmd += c.Copy(c.varPath(b), c.argPath("internal/string/compare", "b"))
	cmd += c.SetArg("internal/string/compare", "res", nbt.NewString(res))
	cmd += c.SetArg("internal/string/compare", "storage", nbt.NewString(c.storage))
	cmd += c.Call("mcb:internal/string/compare")
	return
}
