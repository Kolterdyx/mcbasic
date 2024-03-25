package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"strings"
)

func (c *Compiler) cmd(format string, a ...interface{}) string {
	str := fmt.Sprintf(format, a...)
	// if $ is found anywhere in the string, it is a macro
	if strings.Contains(str, "$") {
		return "$" + str
	}
	return str
}

func (c *Compiler) add(var1 string, var2 string) string {
	cmd := c.declareVar("res")
	cmd += c.cmd("scoreboard players operation %s %s += %s %s\n", "res", c.Namespace, var1, c.Namespace)
	cmd += c.cmd("scoreboard players operation %s %s += %s %s\n", "res", c.Namespace, var2, c.Namespace)
	c.resStore = ScoreRes
	return cmd
}

func (c *Compiler) sub(var1 string, var2 string) string {
	cmd := c.declareVar("res")
	cmd += c.cmd("scoreboard players operation %s %s += %s %s\n", "res", c.Namespace, var1, c.Namespace)
	cmd += c.cmd("scoreboard players operation %s %s -= %s %s\n", "res", c.Namespace, var2, c.Namespace)
	c.resStore = ScoreRes
	return cmd
}

func (c *Compiler) mul(var1 string, var2 string) string {
	cmd := c.declareVar("res")
	cmd += c.cmd("scoreboard players operation %s %s += %s %s\n", "res", c.Namespace, var1, c.Namespace)
	cmd += c.cmd("scoreboard players operation %s %s *= %s %s\n", "res", c.Namespace, var2, c.Namespace)
	cmd += c.declareVarWithVal("res1", "1000")
	cmd += c.cmd("scoreboard players operation %s %s /= %s %s\n", "res", c.Namespace, "res1", c.Namespace)
	c.resStore = ScoreRes
	return cmd
}

func (c *Compiler) div(var1 string, var2 string) string {
	cmd := c.declareVar("res")
	cmd += c.cmd("scoreboard players operation %s %s += %s %s\n", "res", c.Namespace, var1, c.Namespace)
	cmd += c.declareVarWithVal("res1", "1000")
	cmd += c.cmd("scoreboard players operation %s %s *= %s %s\n", "res", c.Namespace, "res1", c.Namespace)
	cmd += c.cmd("scoreboard players operation %s %s /= %s %s\n", "res", c.Namespace, var2, c.Namespace)
	c.resStore = ScoreRes
	return cmd
}

func (c *Compiler) mod(var1 string, var2 string) string {
	cmd := c.declareVar("res")
	cmd += c.cmd("scoreboard players operation %s %s += %s %s\n", "res", c.Namespace, var1, c.Namespace)
	cmd += c.cmd("scoreboard players operation %s %s %%= %s %s\n", "res", c.Namespace, var2, c.Namespace)
	c.resStore = ScoreRes
	return cmd
}

func (c *Compiler) comp(op string, var1 string, var2 string) string {
	cmd := c.declareVar("res")
	switch op {
	case "<":
		fallthrough
	case "<=":
		fallthrough
	case ">":
		fallthrough
	case ">=":
		cmd += c.cmd("execute if score %s %s %s %s %s run scoreboard players set %s %s 1\n", var1, c.Namespace, op, var2, c.Namespace, "res", c.Namespace)
	case "==":
		cmd += c.cmd("execute if score %s %s = %s %s run scoreboard players set %s %s 1\n", var1, c.Namespace, var2, c.Namespace, "res", c.Namespace)
	case "!=":
		cmd += c.cmd("execute unless score %s %s = %s %s run scoreboard players set %s %s 1\n", var1, c.Namespace, var2, c.Namespace, "res", c.Namespace)
	default:
		panic("Unknown comparison operator")
	}
	c.resStore = ScoreRes
	return cmd
}

func (c *Compiler) and(var1 string, var2 string) string {
	cmd := c.declareVar("res")
	cmd += c.cmd("execute unless score %s %s matches 0 unless score %s %s matches 0 run scoreboard players set %s %s 1\n", var1, c.Namespace, var2, c.Namespace, "res", c.Namespace)
	c.resStore = ScoreRes
	return cmd
}

func (c *Compiler) or(var1 string, var2 string) string {
	cmd := c.declareVar("res")
	cmd += c.cmd("execute if score %s %s matches 1 if score %s %s matches 1 run scoreboard players set %s %s 1\n", var1, c.Namespace, var2, c.Namespace, "res", c.Namespace)
	c.resStore = ScoreRes
	return cmd
}

func (c *Compiler) not(var1 string) string {
	cmd := c.declareVar("res")
	cmd += c.comp("==", var1, "0")
	return cmd
}

func (c *Compiler) assignExpr(var1 string, var2 string) string {
	if c.resStore == ScoreRes {
		return c.cmd("scoreboard players operation %s %s = %s %s\n", var1, c.Namespace, var2, c.Namespace)
	} else if c.resStore == DataRes {
		return c.cmd("data modify storage %s:vars %s set from storage %s:vars %s\n", c.Namespace, var1, c.Namespace, var2)
	}
	return ""
}

func (c *Compiler) assignVal(var1 string, var2 string) string {
	if c.resStore == ScoreRes {
		return c.cmd("scoreboard players set %s %s %s\n", var1, c.Namespace, var2)
	} else if c.resStore == DataRes {
		return c.cmd("data modify storage %s:vars %s set value %s\n", c.Namespace, var1, var2)
	}
	return ""
}

func (c *Compiler) print() string {
	return c.call("__print__")
}

func (c *Compiler) call(funcName string) string {
	return c.cmd("function %s:%s with storage %s:func_args.%s\n", c.Namespace, funcName, c.Namespace, funcName)
}

func (c *Compiler) declareVar(varName string) string {
	if varName == "res" {
		c.resStore = ScoreRes
	}
	return c.cmd("scoreboard players set %s %s 0\n", varName, c.Namespace)
}

func (c *Compiler) declareVarWithVarVal(varName string, val string) string {
	cmd := c.declareVar(varName)
	cmd += c.assignExpr(varName, val)
	if varName == "res" {
		c.resStore = ScoreRes
	}
	return cmd
}

func (c *Compiler) declareVarWithVal(varName string, val string) string {
	cmd := c.declareVar(varName)
	cmd += c.assignVal(varName, val)
	if varName == "res" {
		c.resStore = ScoreRes
	}
	return cmd
}

func (c *Compiler) declareVarWithStringVal(varName string, val string) string {
	cmd := c.declareVar(varName)
	val = strings.ReplaceAll(val, "\"", "\\\"")
	cmd += c.cmd("data modify storage %s:vars %s set value \"%s\"\n", c.Namespace, varName, val)
	if varName == "res" {
		c.resStore = DataRes
	}
	return cmd
}

type Arg struct {
	Function string
	Locator  string
	Value    string
}

func (c *Compiler) storeArg(arg Arg) string {
	return c.cmd("data modify storage %s:func_args.%s %s set value %s\n", c.Namespace, arg.Function, arg.Locator, arg.Value)
}

func (c *Compiler) storeArgs(args ...Arg) string {
	cmd := ""
	for _, arg := range args {
		cmd += c.storeArg(arg)
	}
	return cmd
}

func (c *Compiler) storeArgFromVar(function string, arg string, varName string, location ResStore) string {
	if location == ScoreRes {
		return c.cmd("execute store result storage %s:func_args.%s %s int 1 run scoreboard players get %s %s\n", c.Namespace, function, arg, varName, c.Namespace)
	} else if location == DataRes {
		return c.cmd("data modify storage %s:func_args.%s %s set from storage %s:vars %s\n", c.Namespace, function, arg, c.Namespace, varName)
	}
	return ""
}

func (c *Compiler) macro(varName string) string {
	var varType tokens.TokenType
	for i, p := range c.functionArgs[c.currentFunction.Name.Lexeme] {
		if p == varName {
			varType = c.currentFunction.Types[i].Type
		}
	}
	if varType == tokens.NumberType {
		c.resStore = ScoreRes
		return fmt.Sprintf("$(%s)", varName)
	} else if varType == tokens.StringType {
		c.resStore = DataRes
		return fmt.Sprintf(`"$(%s)"`, varName)
	}
	return ""

}
