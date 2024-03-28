package ops

import (
	"fmt"
	"strconv"
)

// Sc stands for "scoreboard". It moves a variable from data storage to the scoreboard for further processing.
func (o *Op) Sc(varName string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:%s %s\n", varName, o.Namespace, o.Namespace, VarPath, varName)
}

// St stands for "store". It moves a variable from the scoreboard to data storage for further processing.
func (o *Op) St(varName string) string {
	return fmt.Sprintf("execute store result storage %s:%s %s int 1 run scoreboard players get %s %s\n", o.Namespace, VarPath, varName, varName, varName)
}

// Set sets a variable to a specific value.
func (o *Op) Set(varName string, literalValue string) string {
	// if literalValue is a number, store as is, otherwise, wrap in quotes
	if _, err := strconv.Atoi(literalValue); err != nil {
		literalValue = fmt.Sprintf("\"%s\"", literalValue)
	}
	return fmt.Sprintf("data modify storage %s:%s %s set value %s\n", o.Namespace, VarPath, varName, literalValue)
}

// SetMacro sets a variable to a macro value.
func (o *Op) SetMacro(varName string, macro string) string {
	return fmt.Sprintf("data modify storage %s:%s %s set value %s\n", o.Namespace, VarPath, varName, macro)
}

// Cp stands for "copy". It copies the value of one variable to another.
func (o *Op) Cp(fromVar string, toVar string) string {
	return fmt.Sprintf("data modify storage %s:%s %s set from storage %s:%s %s\n", o.Namespace, VarPath, toVar, o.Namespace, VarPath, fromVar)
}
