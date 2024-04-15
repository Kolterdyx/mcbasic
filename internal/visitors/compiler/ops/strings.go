package ops

func (o *Op) Concat(var1 string, var2 string, result string) string {
	cmd := ""
	cmd += o.LoadArgConst("internal/concat", "res", result)
	cmd += o.CallFunction("internal/concat", map[string]string{"a": var1, "b": var2}, RX)
	return cmd
}
