package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
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
		c.opHandler.DoubleSqrt("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)

	// Trigonometric functions
	c.createFunction(
		"math:cos",
		c.opHandler.DoubleCos("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:sin",
		c.opHandler.DoubleSin("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:tan",
		c.opHandler.DoubleTan("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:acos",
		c.opHandler.DoubleAcos("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:asin",
		c.opHandler.DoubleAsin("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:atan",
		c.opHandler.DoubleAtan("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)

	// Rounding functions
	c.createFunction(
		"math:floor",
		c.opHandler.DoubleFloor("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:ceil",
		c.opHandler.DoubleCeil("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
	c.createFunction(
		"math:round",
		c.opHandler.DoubleRound("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: types.DoubleType},
		},
		types.DoubleType,
	)
}

func (c *Compiler) baseFunctions() {

	cleanCall := ""
	if c.Config.Project.CleanBeforeInit {
		cleanCall = c.opHandler.Call("mcb:internal/clean", "")
	}

	c.createFunction(
		"mcb:internal/init",
		fmt.Sprintf("scoreboard objectives add %s dummy\n", c.Namespace)+
			c.opHandler.MakeConst("0", ops.CALL, false)+
			c.opHandler.MoveScore(ops.CALL, ops.CALL)+
			cleanCall+
			c.opHandler.LoadArgConst("log", "text", "MCB pack loaded", true)+
			c.opHandler.Call("mcb:log", "")+
			c.opHandler.Call("internal/struct_definitions", "")+
			c.opHandler.Call("main", ""),
		[]interfaces.FuncArg{},
		types.VoidType,
	)
	c.createFunction(
		"mcb:internal/clean",
		"data remove storage example:data vars\n"+
			"data remove storage example:data structs\n"+
			"data remove storage example:data args",
		[]interfaces.FuncArg{},
		types.VoidType,
	)
	c.createFunction(
		"internal/tick",
		c.opHandler.Call("tick", "")+
			c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 1000..", ops.CALL, c.Namespace), true, c.opHandler.MakeConst("0", ops.CALL)),
		[]interfaces.FuncArg{},
		types.VoidType,
	)
}
