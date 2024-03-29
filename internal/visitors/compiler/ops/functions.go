package ops

import "fmt"

// ArgLoad loads an argument from a variable.
func (o *Op) ArgLoad(funcName string, argName string, varName string) string {
	return fmt.Sprintf("data modify storage %s:%s.%s %s set from storage %s:%s.%s %s\n", o.Namespace, ArgPath, funcName, argName, o.Namespace, VarPath, o.Scope, varName)
}

// Call calls a function.
func (o *Op) Call(funcName string) string {
	return fmt.Sprintf("function %s:%s with storage %s:%s.%s\n", o.Namespace, funcName, o.Namespace, ArgPath, funcName)
}
