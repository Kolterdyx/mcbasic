package ir

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

func (c *Compiler) compileBuiltins() []Function {
	ir := c.math()
	ir = append(ir, c.baseFunctions()...)
	return ir
}

func (c *Compiler) math() []Function {
	return []Function{
		c.createIRFunction(
			"math:sqrt",
			c.DoubleSqrt("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),

		// Trigonometric functions
		c.createIRFunction(
			"math:cos",
			c.DoubleCos("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.createIRFunction(
			"math:sin",
			c.DoubleSin("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.createIRFunction(
			"math:tan",
			c.DoubleTan("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.createIRFunction(
			"math:acos",
			c.Copy(
				fmt.Sprintf("%s.acos.x", ArgPath),
				fmt.Sprintf("%s.x", VarPath),
			)+
				c.DoubleAcos("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.createIRFunction(
			"math:asin",
			c.Copy(
				fmt.Sprintf("%s.asin.x", ArgPath),
				fmt.Sprintf("%s.x", VarPath),
			)+
				c.DoubleAsin("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.createIRFunction(
			"math:atan",
			c.Copy(
				fmt.Sprintf("%s.atan.x", ArgPath),
				fmt.Sprintf("%s.x", VarPath),
			)+
				c.DoubleAtan("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),

		// Rounding functions
		c.createIRFunction(
			"math:floor",
			c.DoubleFloor("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.createIRFunction(
			"math:ceil",
			c.DoubleCeil("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.createIRFunction(
			"math:round",
			c.DoubleRound("x", RET)+
				c.Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
	}
}

func (c *Compiler) baseFunctions() []Function {

	funcs := make([]Function, 0)

	cleanCall := ""
	if c.Config.CleanBeforeInit {
		cleanCall = c.Call("mcb:internal/clean")
	}

	initSource := ""
	initSource += cleanCall
	initSource += fmt.Sprintf("scoreboard objectives add %s dummy\n", c.Namespace)
	if c.Config.Debug {
		initSource += fmt.Sprintf("scoreboard objectives setdisplay sidebar %s\n", c.Namespace)
	}
	initSource += c.SetVar(CALL, nbt.NewInt(0))
	initSource += c.Load(CALL, CALL)
	initSource += c.SetArg("mcb:log", "text", nbt.NewString("MCB pack loaded"))
	initSource += c.Call("mcb:log")
	initSource += c.Call("internal/struct_definitions")
	initSource += c.Call("main")
	initSource += c.Ret()
	funcs = append(
		funcs,
		c.createIRFunction(
			"mcb:internal/init",
			initSource,
			[]interfaces.TypedIdentifier{},
			types.VoidType,
		),
		c.createIRFunction(
			"mcb:internal/clean",
			//fmt.Sprintf("scoreboard objectives remove %s\n", c.Namespace)+
			c.Remove(VarPath)+
				c.Remove(StructPath)+
				c.Remove(ArgPath)+
				c.Ret(),
			[]interfaces.TypedIdentifier{},
			types.VoidType,
		),
	)
	tickSource := ""
	tickSource += c.SetVar(CALL, nbt.NewInt(0))
	tickSource += c.Load(CALL, CALL)

	maxCallCounterReg := c.makeReg(CALL)
	compResult := c.makeReg(RET)
	funcs = append(
		funcs,
		c.createIRFunction(
			"internal/tick",
			c.Call("tick")+
				c.Set(maxCallCounterReg, nbt.NewInt(MaxCallCounter))+
				c.IntCompare(CALL, maxCallCounterReg, tokens.GreaterEqual, compResult)+
				c.If(compResult, tickSource)+
				c.Ret(),
			[]interfaces.TypedIdentifier{},
			types.VoidType,
		),
		c.createIRFunction(
			"tick",
			c.Ret(),
			[]interfaces.TypedIdentifier{},
			types.VoidType,
		),
	)
	return funcs
}
