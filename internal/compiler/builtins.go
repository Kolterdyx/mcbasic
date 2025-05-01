package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/compiler/mapping"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

func (c *Compiler) createBuiltinFunctions() {
	c.math()
	c.baseFunctions()
}

func (c *Compiler) math() {
	// Others
	c.createFunction(
		"math:sqrt",
		c.commandMapper.DoubleSqrt("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)

	// Trigonometric functions
	c.createFunction(
		"math:cos",
		c.commandMapper.DoubleCos("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:sin",
		c.commandMapper.DoubleSin("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:tan",
		c.commandMapper.DoubleTan("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:acos",
		c.commandMapper.MoveRaw(
			fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf("%s.acos.x", mapping.ArgPath),
			fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf(c.commandMapper.Cs("%s.x"), mapping.VarPath),
		)+
			c.commandMapper.DoubleAcos("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:asin",
		c.commandMapper.MoveRaw(
			fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf("%s.asin.x", mapping.ArgPath),
			fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf(c.commandMapper.Cs("%s.x"), mapping.VarPath),
		)+
			c.commandMapper.DoubleAsin("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:atan",
		c.commandMapper.MoveRaw(
			fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf("%s.atan.x", mapping.ArgPath),
			fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf(c.commandMapper.Cs("%s.x"), mapping.VarPath),
		)+
			c.commandMapper.DoubleAtan("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)

	// Rounding functions
	c.createFunction(
		"math:floor",
		c.commandMapper.DoubleFloor("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:ceil",
		c.commandMapper.DoubleCeil("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:round",
		c.commandMapper.DoubleRound("x", mapping.RET)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
}

func (c *Compiler) baseFunctions() {

	cleanCall := ""
	if c.Config.CleanBeforeInit {
		cleanCall = c.commandMapper.Call("mcb:internal/clean", "")
	}

	initSource := ""
	initSource += cleanCall
	initSource += fmt.Sprintf("scoreboard objectives add %s dummy\n", c.Namespace)
	if c.Config.Debug {
		initSource += fmt.Sprintf("scoreboard objectives setdisplay sidebar %s\n", c.Namespace)
	}
	initSource += c.commandMapper.MakeConst(nbt.NewInt(0), mapping.CALL)
	initSource += c.commandMapper.MoveScore(mapping.CALL, mapping.CALL)
	initSource += c.commandMapper.LoadArgConst("log", "text", nbt.NewString("MCB pack loaded"))
	initSource += c.commandMapper.Call("mcb:log", "")
	initSource += c.commandMapper.Call("internal/struct_definitions", "")
	initSource += c.commandMapper.Call("main", "")
	initSource += c.commandMapper.Return()
	c.createFunction(
		"mcb:internal/init",
		initSource,
		[]interfaces.FuncArg{},
		types.VoidType,
	)
	c.createFunction(
		"mcb:internal/clean",
		fmt.Sprintf("scoreboard objectives remove %s\n", c.Namespace)+
			fmt.Sprintf("data remove storage %s:data vars\n", c.Namespace)+
			fmt.Sprintf("data remove storage %s:data structs\n", c.Namespace)+
			fmt.Sprintf("data remove storage %s:data args\n", c.Namespace)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{},
		types.VoidType,
	)
	tickSource := ""
	tickSource += c.commandMapper.MakeConst(nbt.NewInt(0), mapping.CALL)
	tickSource += c.commandMapper.MoveScore(mapping.CALL, mapping.CALL)
	c.createFunction(
		"internal/tick",
		c.commandMapper.Call("tick", "")+
			c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches %d..", mapping.CALL, c.Namespace, mapping.MaxCallCounter), true, tickSource)+
			c.commandMapper.Return(),
		[]interfaces.FuncArg{},
		types.VoidType,
	)
	c.createFunction(
		"tick",
		c.commandMapper.Return(),
		[]interfaces.FuncArg{},
		types.VoidType,
	)
}
