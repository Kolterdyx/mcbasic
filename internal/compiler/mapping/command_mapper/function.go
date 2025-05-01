package command_mapper

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (c *CommandMapper) Call(funcName string, res string) string {
	data := nbt.NewCompound()
	ns, fn := c.baseFuncName(funcName)
	data.Set("function", nbt.NewString(fn))
	data.Set("function_namespace", nbt.NewString(ns))
	data.Set("namespace", nbt.NewString(c.Namespace))
	data.Set("args", nbt.NewString(fmt.Sprintf("%s:data %s.%s", c.Namespace, c.ArgPath, fn)))
	if res == "" {
		res = c.RET
	}
	data.Set("ret", nbt.NewString(res))
	return fmt.Sprintf("function mcb:internal/call %s\n", data.ToString())
}

func (c *CommandMapper) LoadArgs(funcName string, args map[string]string) string {
	cmd := ""
	for k, v := range args {
		cmd += c.LoadArg(funcName, k, v)
	}
	return cmd
}

func (c *CommandMapper) LoadArg(funcName, argName string, varName string) string {
	_, fn := c.baseFuncName(funcName)
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set from storage %s:data %s.%s\n", c.Namespace, c.ArgPath, fn, argName, c.Namespace, c.VarPath, c.Cs(varName))
}

func (c *CommandMapper) LoadArgRaw(funcName, argName string, varName string) string {
	_, fn := c.baseFuncName(funcName)
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set from storage %s:data %s\n", c.Namespace, c.ArgPath, fn, argName, c.Namespace, varName)
}

func (c *CommandMapper) LoadArgConst(funcName, argName string, value nbt.Value) string {
	_, fn := c.baseFuncName(funcName)
	return fmt.Sprintf("data modify storage %s:data %s.%s.%s set value %s\n", c.Namespace, c.ArgPath, fn, argName, value.ToString())
}

// baseFuncName returns the base function name and the namespace:
// namespace, funcName
func (c *CommandMapper) baseFuncName(funcName string) (string, string) {
	if strings.Contains(funcName, ":") {
		split := strings.Split(funcName, ":")
		if len(split) == 2 {
			return split[0], split[1]
		}
		log.Fatal("Function name is not valid: ", funcName)
	}
	return c.Namespace, funcName
}
