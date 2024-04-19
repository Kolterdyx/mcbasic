package ops

import (
	"fmt"
	"strconv"
	"strings"
)

func (o *Op) Move(from string, to string) string {
	return fmt.Sprintf("data modify storage %s:%s %s set from storage %s:%s %s\n", o.Namespace, VarPath, to, o.Namespace, VarPath, from)
}

func (o *Op) MoveConst(value string, to string) string {
	if _, err := strconv.Atoi(value); err != nil && !(value[0] == '$' && value[1] == '(' && value[len(value)-1] == ')') && !(value[0] == '"' && value[len(value)-1] == '"') {
		value = strconv.Quote(value)
	}
	return fmt.Sprintf("data modify storage %s:%s %s set value %s\n", o.Namespace, VarPath, to, value)
}

func (o *Op) MoveFixedConst(value string, to string) string {
	vw := to + ".whole"
	vf := to + ".fract"
	values := strings.Split(value, ".")
	if len(values) != 2 {
		values = append(values, "0")
	}
	fmt.Printf("values: %v\n", values)
	cmd := ""
	w, _ := strconv.ParseInt(values[0], 10, 64)
	f, _ := strconv.ParseInt(values[1], 10, 64)
	cmd += o.MoveConst(strconv.FormatInt(w, 10), vw)
	cmd += o.MoveConst(strconv.FormatInt(f, 10), vf)
	return cmd
}

func (o *Op) MoveScore(from string, to string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:%s %s\n", to, o.Namespace, o.Namespace, VarPath, from)
}

func (o *Op) LoadScore(from string, to string) string {
	return fmt.Sprintf("execute store result storage %s:%s %s int 1 run scoreboard players get %s %s\n", o.Namespace, VarPath, to, from, o.Namespace)
}

func (o *Op) Inc(varName string) string {
	cmd := ""
	cmd += o.MoveScore(varName, varName)
	cmd += fmt.Sprintf("scoreboard players add %s %s 1\n", varName, o.Namespace)
	cmd += o.LoadScore(varName, varName)
	return cmd
}
