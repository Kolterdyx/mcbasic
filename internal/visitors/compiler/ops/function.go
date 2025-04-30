package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
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
	data := nbt.NewCompound()
	ns, fn := o.baseFuncName(funcName)
	data.Set("function", nbt.NewString(fn))
	data.Set("namespace", nbt.NewString(ns))
	data.Set("args", nbt.NewString(fmt.Sprintf("%s:data %s.%s", o.Namespace, ArgPath, fn)))
	if res == "" {
		res = RET
	}
	data.Set("ret", nbt.NewString(res))
	return fmt.Sprintf("function mcb:internal/call %s\n", data.ToString())
}

func (o *Op) LoadArgs(funcName string, args map[string]string) string {
	cmd := ""
	for k, v := range args {
		cmd += o.LoadArg(funcName, k, v)
	}
	return cmd
}

func (o *Op) LoadArg(funcName, argName string, varName string) string {
	_, fn := o.baseFuncName(funcName)
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set from storage %s:data %s.%s\n", o.Namespace, ArgPath, fn, argName, o.Namespace, VarPath, Cs(varName))
}

func (o *Op) LoadArgRaw(funcName, argName string, varName string) string {
	_, fn := o.baseFuncName(funcName)
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set from storage %s:data %s\n", o.Namespace, ArgPath, fn, argName, o.Namespace, varName)
}

func (o *Op) LoadArgConst(funcName, argName string, value string, quote ...bool) string {
	// if value is not numeric, wrap it in quotes
	if len(quote) > 0 && quote[0] {
		value = strconv.Quote(value)
	}
	_, fn := o.baseFuncName(funcName)
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set value %s\n", o.Namespace, ArgPath, fn, argName, value)
}

// baseFuncName returns the base function name and the namespace:
// namespace, funcName
func (o *Op) baseFuncName(funcName string) (string, string) {
	if strings.Contains(funcName, ":") {
		split := strings.Split(funcName, ":")
		if len(split) == 2 {
			return split[0], split[1]
		}
		log.Fatal("Function name is not valid: ", funcName)
	}
	return o.Namespace, funcName
}
