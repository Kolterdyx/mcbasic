package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
)

func (o *Op) Concat(var1, var2, result string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/concat", "res", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, RET)))
	cmd += o.LoadArgConst("internal/concat", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArg("internal/concat", "a", var1)
	cmd += o.LoadArg("internal/concat", "b", var2)
	cmd += o.Call("mcb:internal/concat", result)
	return cmd
}

func (o *Op) Size(var1, result string) string {
	return fmt.Sprintf("execute store result storage %s:data %s.%s int 1 run data get storage %s:data %s.%s\n", o.Namespace, VarPath, result, o.Namespace, VarPath, var1)
}

func (o *Op) SliceString(from, start, end, result string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/slice", "res", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, RET)))
	cmd += o.LoadArgConst("internal/slice", "storage", nbt.NewString(fmt.Sprintf("%s:data", o.Namespace)))
	cmd += o.LoadArgConst("internal/slice", "from", nbt.NewString(fmt.Sprintf("%s.%s", VarPath, from)))
	cmd += o.LoadArg("internal/slice", "start", start)
	cmd += o.LoadArg("internal/slice", "end", end)
	cmd += o.Call("mcb:internal/slice", result)
	return cmd
}
