package ops

import "fmt"

// RegWrite stands for "register write". It moves a variable from data storage to a register for further processing.
func (o *Op) RegWrite(varName string, regName string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:%s.%s %s\n", regName, o.Namespace, o.Namespace, VarPath, o.Scope, varName)
}

// RegLoad stands for "register load". It writes a value to a register.
func (o *Op) RegLoad(value string, regName string) string {
	return fmt.Sprintf("scoreboard players set %s %s %s\n", regName, o.Namespace, value)
}

// RegShift stands for "register shift". It moves a variable from one register to another.
func (o *Op) RegShift(regFrom string, regTo string) string {
	return fmt.Sprintf("scoreboard players operation %s %s = %s %s\n", regTo, o.Namespace, regFrom, o.Namespace)
}

// RegSave stands for "register save". It moves a variable from a register to data storage for further processing.
func (o *Op) RegSave(varName string, regName string) string {
	return fmt.Sprintf("execute store result storage %s:%s.%s %s int 1 run scoreboard players get %s %s\n", o.Namespace, VarPath, o.Scope, varName, regName, o.Namespace)
}
