package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
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
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)

	// Trigonometric functions
	c.createFunction(
		"math:cos",
		c.opHandler.DoubleCos("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)
	c.createFunction(
		"math:sin",
		c.opHandler.DoubleSin("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)
	c.createFunction(
		"math:tan",
		c.opHandler.DoubleTan("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)
	c.createFunction(
		"math:acos",
		c.opHandler.DoubleAcos("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)
	c.createFunction(
		"math:asin",
		c.opHandler.DoubleAsin("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)
	c.createFunction(
		"math:atan",
		c.opHandler.DoubleAtan("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)

	// Rounding functions
	c.createFunction(
		"math:floor",
		c.opHandler.DoubleFloor("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)
	c.createFunction(
		"math:ceil",
		c.opHandler.DoubleCeil("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)
	c.createFunction(
		"math:round",
		c.opHandler.DoubleRound("x", ops.RET),
		[]interfaces.FuncArg{
			{Name: "x", Type: expressions.DoubleType},
		},
		expressions.DoubleType,
	)
}

func (c *Compiler) baseFunctions() {
	c.createFunction(
		"mcb:internal/init",
		fmt.Sprintf("scoreboard objectives add %s dummy\n", c.Namespace)+
			c.opHandler.MoveConst("0", ops.CALL)+
			c.opHandler.MoveScore(ops.CALL, ops.CALL)+
			c.opHandler.LoadArgConst("log", "text", "MCB pack loaded")+
			c.opHandler.Call("mcb:log", "")+
			c.opHandler.Call("main", ""),
		[]interfaces.FuncArg{},
		expressions.VoidType,
	)
	c.createFunction(
		"internal/tick",
		c.opHandler.Call("tick", "")+
			c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 1000..", ops.CALL, c.Namespace), true, c.opHandler.MoveConst("0", ops.CALL)),
		[]interfaces.FuncArg{},
		expressions.VoidType,
	)
}
