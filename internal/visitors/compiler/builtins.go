package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
)

func (c *Compiler) createBuiltinFunctions() {
	c.math()
	c.baseFunctions()
}

func (c *Compiler) math() {
	// Others
	c.createFunction(
		"math:sqrt",
		c.opHandler.DoubleSqrt("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)

	// Trigonometric functions
	c.createFunction(
		"math:cos",
		c.opHandler.DoubleCos("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:sin",
		c.opHandler.DoubleSin("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:tan",
		c.opHandler.DoubleTan("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:acos",
		c.opHandler.DoubleAcos("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:asin",
		c.opHandler.DoubleAsin("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:atan",
		c.opHandler.DoubleAtan("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)

	// Rounding functions
	c.createFunction(
		"math:floor",
		c.opHandler.DoubleFloor("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:ceil",
		c.opHandler.DoubleCeil("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:round",
		c.opHandler.DoubleRound("x", ops.RET)+
			c.opHandler.Return(),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
}

func (c *Compiler) baseFunctions() {

	cleanCall := ""
	if c.Config.CleanBeforeInit {
		cleanCall = c.opHandler.Call("mcb:internal/clean", "")
	}

	initSource := ""
	initSource += cleanCall
	initSource += fmt.Sprintf("scoreboard objectives add %s dummy\n", c.Namespace)
	if c.Config.Debug {
		initSource += fmt.Sprintf("scoreboard objectives setdisplay sidebar %s\n", c.Namespace)
	}
	initSource += c.opHandler.MakeConst(nbt.NewInt(0), ops.CALL)
	initSource += c.opHandler.MoveScore(ops.CALL, ops.CALL)
	initSource += c.opHandler.LoadArgConst("log", "text", nbt.NewString("MCB pack loaded"))
	initSource += c.opHandler.Call("mcb:log", "")
	initSource += c.opHandler.Call("internal/struct_definitions", "")
	initSource += c.opHandler.Call("main", "")
	initSource += c.opHandler.Return()
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
			c.opHandler.Return(),
		[]interfaces.FuncArg{},
		types.VoidType,
	)
	c.createFunction(
		"internal/tick",
		c.opHandler.Call("tick", "")+
			c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 1000..", ops.CALL, c.Namespace), true, c.opHandler.MakeConst(nbt.NewInt(0), ops.CALL))+
			c.opHandler.Return(),
		[]interfaces.FuncArg{},
		types.VoidType,
	)
	c.createFunction(
		"tick",
		c.opHandler.Return(),
		[]interfaces.FuncArg{},
		types.VoidType,
	)
}
