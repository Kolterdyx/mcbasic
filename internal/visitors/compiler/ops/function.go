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
		cmd += fmt.Sprintf("function mcb:internal/call {function:'%s',args:'%s:data %s.%s',namespace:'%s'}\n", funcName, o.Namespace, ArgPath, o.baseFuncName(funcName), o.Namespace)
	} else {
		cmd += fmt.Sprintf("function mcb:internal/call {function:'%s:%s',args:'%s:data %s.%s',namespace:'%s'}\n", o.Namespace, funcName, o.Namespace, ArgPath, funcName, o.Namespace)
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

func (o *Op) LoadArgRaw(funcName, argName string, varName string) string {
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set from storage %s:data %s\n", o.Namespace, ArgPath, o.baseFuncName(funcName), argName, o.Namespace, varName)
}

func (o *Op) LoadArgConst(funcName, argName string, value string, quote ...bool) string {
	// if value is not numeric, wrap it in quotes
	if len(quote) > 0 && quote[0] {
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
