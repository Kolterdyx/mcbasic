package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
)

func (c *Compiler) createBuiltinFunctions() {
	c.utils()
	c.math()
	c.strings()
	c.baseFunctions()
}

func (c *Compiler) utils() {
	c.createFunction(
		"mcb:print",
		`$tellraw @a {text:'$(text)'}`,
		[]interfaces.FuncArg{
			{Name: "text", Type: expressions.StringType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:log",
		`$tellraw @a[tag=mcblog] {text:'$(text)',color:gray,italic:true}`,
		[]interfaces.FuncArg{
			{Name: "text", Type: expressions.StringType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:exec",
		`$execute run $(command)`,
		[]interfaces.FuncArg{
			{Name: "command", Type: expressions.StringType},
		},
		expressions.VoidType,
	)
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

func (c *Compiler) strings() {
	c.createFunction(
		"mcb:internal/concat",
		`$data modify storage $(storage) $(res) set value "$(a)$(b)"`,
		[]interfaces.FuncArg{
			{Name: "storage", Type: expressions.StringType},
			{Name: "res", Type: expressions.StringType},
			{Name: "a", Type: expressions.StringType},
			{Name: "b", Type: expressions.VoidType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:internal/slice",
		`$data modify storage $(storage) $(res) set string storage $(storage) $(from) $(start) $(end)`,
		[]interfaces.FuncArg{
			{Name: "storage", Type: expressions.StringType},
			{Name: "res", Type: expressions.StringType},
			{Name: "from", Type: expressions.StringType},
			{Name: "start", Type: expressions.IntType},
			{Name: "end", Type: expressions.IntType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:len",
		fmt.Sprintf("$data modify storage %s:%s %s set value \"$(from)\"\n", c.Namespace, ops.VarPath, ops.RET)+
			fmt.Sprintf("execute store result storage %s:%s %s int 1 run data get storage %s:%s %s\n", c.Namespace, ops.VarPath, ops.RET, c.Namespace, ops.VarPath, ops.RET),
		[]interfaces.FuncArg{
			{Name: "from", Type: expressions.StringType},
		},
		expressions.IntType,
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
		c.opHandler.Call("tick", ""),
		[]interfaces.FuncArg{},
		expressions.VoidType,
	)
}
