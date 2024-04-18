package ops

import "fmt"

func (o *Op) Concat(var1 string, var2 string, result string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/concat", "res", result)
	cmd += o.CallFunction("internal/concat", map[string]string{"a": var1, "b": var2}, RX)
	return cmd
}

func (o *Op) LenString(var1 string, result string) string {
	return fmt.Sprintf("execute store result storage %s:%s %s int 1 run data get storage %s:%s %s\n", o.Namespace, VarPath, result, o.Namespace, VarPath, var1)
}
