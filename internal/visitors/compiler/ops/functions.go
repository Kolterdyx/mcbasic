package ops

import "fmt"

// ArgLoad loads an argument from a variable.
func (o *Op) ArgLoad(funcName string, argName string, varName string) string {
	return fmt.Sprintf("data modify storage %s:%s.%s %s set from storage %s:%s.%s %s\n", o.Namespace, ArgPath, funcName, argName, o.Namespace, VarPath, o.Scope, cs(varName))
}

// Call calls a function.
func (o *Op) Call(funcName string) string {
	cmd := ""
	cmd += o.Inc(CALL)
	cmd += o.RegSave(RX, CALL)
	cmd += o.ArgLoad(funcName, "__call__", RX)
	cmd += fmt.Sprintf("execute store result storage %s:%s.%s %s int 1 run function %s:%s with storage %s:%s.%s\n", o.Namespace, VarPath, o.Scope, cs(RX), o.Namespace, funcName, o.Namespace, ArgPath, funcName)
	return cmd
}
