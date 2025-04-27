package ops

import (
	"fmt"
	"strconv"
)

func (o *Op) MoveRaw(storageFrom, pathFrom, storageTo, pathTo string) string {
	return fmt.Sprintf("data modify storage %s %s set from storage %s %s\n", storageTo, pathTo, storageFrom, pathFrom)
}

func (o *Op) Move(from, to string) string {
	return o.MoveRaw(
		fmt.Sprintf("%s:data", o.Namespace), fmt.Sprintf("%s.%s", VarPath, from),
		fmt.Sprintf("%s:data", o.Namespace), fmt.Sprintf("%s.%s", VarPath, to),
	)
}

func (o *Op) MakeConst(value, to string, wrapInQuotes ...bool) string {
	wrapInQuote := true
	if len(wrapInQuotes) > 0 {
		wrapInQuote = wrapInQuotes[0]
	}
	if !(value[0] == '$' && value[1] == '(' && value[len(value)-1] == ')') &&
		!(value[0] == '"' && value[len(value)-1] == '"') &&
		wrapInQuote {
		value = strconv.Quote(value)
	}
	// if the value is a float, add a trailing L to store it as a long
	n, err := strconv.ParseFloat(value, 64)
	_, err2 := strconv.ParseInt(value, 10, 64)
	if err == nil && err2 != nil {
		value = fmt.Sprintf("%sd", strconv.FormatFloat(n, 'f', -1, 64))
	}
	return fmt.Sprintf("data modify storage %s:data %s.%s set value %s\n", o.Namespace, VarPath, to, value)
}

func (o *Op) MoveScore(from string, to string) string {
	return fmt.Sprintf("execute store result score %s %s run data get storage %s:data %s.%s\n", to, o.Namespace, o.Namespace, VarPath, from)
}

func (o *Op) LoadScore(from string, to string) string {
	return fmt.Sprintf("execute store result storage %s:data %s.%s int 1 run scoreboard players get %s %s\n", o.Namespace, VarPath, to, from, o.Namespace)
}

func (o *Op) Inc(varName string) string {
	cmd := ""
	cmd += o.MoveScore(varName, varName)
	cmd += fmt.Sprintf("scoreboard players add %s %s 1\n", varName, o.Namespace)
	cmd += o.LoadScore(varName, varName)
	return cmd
}
