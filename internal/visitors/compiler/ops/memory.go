package ops

import (
	"fmt"
)

// Sc stands for "scoreboard". It moves a variable from data storage to the scoreboard for further processing.
func (o *Op) Sc(varName string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:%s %s", varName, o.Namespace, o.Namespace, VarPath, varName)
}

// St stands for "store". It moves a variable from the scoreboard to data storage for further processing.
func (o *Op) St(varName string) string {
	return fmt.Sprintf("execute store result storage %s:%s %s run scoreboard players get %s %s", o.Namespace, VarPath, varName, varName, varName)
}

// Rw stands for "register write". It moves a variable from data storage to a register for further processing.
func (o *Op) Rw(varName string, regName string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:%s %s", regName, o.Namespace, o.Namespace, VarPath, varName)
}

// Rs stands for "register save". It moves a variable from a register to data storage for further processing.
func (o *Op) Rs(varName string, regName string) string {
	return fmt.Sprintf("execute store result storage %s:%s %s run scoreboard players get %s %s", o.Namespace, VarPath, varName, regName, regName)
}

// Set sets a variable to a specific value.
func (o *Op) Set(varName string, literalValue string) string {
	return fmt.Sprintf("execute store result storage %s:%s int 1 run data modify storage %s:%s %s set value %s", o.Namespace, VarPath, o.Namespace, VarPath, varName, literalValue)
}

// Cp stands for "copy". It copies the value of one variable to another.
func (o *Op) Cp(fromVar string, toVar string) string {
	return fmt.Sprintf("data modify storage %s:%s %s set from storage %s:%s %s", o.Namespace, VarPath, toVar, o.Namespace, VarPath, fromVar)
}
