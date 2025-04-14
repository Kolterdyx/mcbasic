package ops

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (o *Op) CallFunction(funcName string, args map[string]string, res string) string {
	cmd := ""
	cmd += o.LoadArgs(funcName, args)
	cmd += o.Call(funcName, res)

	return cmd
}

func (o *Op) Call(funcName string, res string) string {
	cmd := ""
	cmd += o.LoadArg(funcName, "__call__", CALL)
	cmd += o.Inc(CALL)
	if strings.Contains(funcName, ":") {
		cmd += fmt.Sprintf("function %s with storage %s:data %s.%s\n", funcName, o.Namespace, ArgPath, o.baseFuncName(funcName))
	} else {
		cmd += fmt.Sprintf("function %s:%s with storage %s:data %s.%s\n", o.Namespace, funcName, o.Namespace, ArgPath, funcName)
	}
	if res == "" {
		return cmd
	}
	cmd += o.Move(RET, Cs(res))

	return cmd
}

func (o *Op) LoadArgs(funcName string, args map[string]string) string {
	cmd := ""
	for k, v := range args {
		cmd += o.LoadArg(funcName, k, v)
	}
	return cmd
}

func (o *Op) LoadArg(funcName, argName string, varName string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set from storage %s:data %s.%s\n", o.Namespace, ArgPath, o.baseFuncName(funcName), argName, o.Namespace, VarPath, Cs(varName))
}

func (o *Op) LoadArgConst(funcName, argName string, value string) string {
	// if value is not numeric, wrap it in quotes
	if _, err := strconv.Atoi(value); err != nil && !(value[0] == '$' && value[1] == '(' && value[len(value)-1] == ')') {
		value = strconv.Quote(value)
	}
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set value %s\n", o.Namespace, ArgPath, o.baseFuncName(funcName), argName, value)
}

func (o *Op) baseFuncName(funcName string) string {
	if strings.Contains(funcName, ":") {
		split := strings.Split(funcName, ":")
		if len(split) == 2 {
			return split[1]
		}
		log.Fatal("Function name is not valid: ", funcName)
	}
	return funcName
}
