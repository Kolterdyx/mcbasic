package ops

import "fmt"

func (o *Op) Concat(var1 string, var2 string, result string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/concat", "res", RET)
	cmd += o.LoadArgConst("internal/concat", "storage", fmt.Sprintf("%s:%s", o.Namespace, VarPath))
	cmd += o.LoadArg("internal/concat", "a", var1)
	cmd += o.LoadArg("internal/concat", "b", var2)
	cmd += o.Call("mcb:internal/concat", result)
	return cmd
}

func (o *Op) SizeString(var1 string, result string) string {
	return fmt.Sprintf("execute store result storage %s:%s %s int 1 run data get storage %s:%s %s\n", o.Namespace, VarPath, result, o.Namespace, VarPath, var1)
}

func (o *Op) Slice(from string, start string, end string, result string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/slice", "res", RET)
	cmd += o.LoadArgConst("internal/slice", "storage", fmt.Sprintf("%s:%s", o.Namespace, VarPath))
	cmd += o.LoadArgConst("internal/slice", "from", from)
	cmd += o.LoadArg("internal/slice", "start", start)
	cmd += o.LoadArg("internal/slice", "end", end)
	cmd += o.Call("mcb:internal/slice", result)
	return cmd
}
