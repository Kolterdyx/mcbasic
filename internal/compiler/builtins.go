package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/ir"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"path"
)

func (c *Compiler) compileBuiltins() []interfaces.Function {
	code := c.math()
	code = append(code, c.baseFunctions()...)
	return code
}

func (c *Compiler) math() []interfaces.Function {
	return []interfaces.Function{
		c.registerIRFunction(
			"math:sqrt",
			ir.NewCode(c.Namespace, c.storage).
				DoubleSqrt("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),

		// Trigonometric functions
		c.registerIRFunction(
			"math:cos",
			ir.NewCode(c.Namespace, c.storage).
				DoubleCos("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.registerIRFunction(
			"math:sin",
			ir.NewCode(c.Namespace, c.storage).
				DoubleSin("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.registerIRFunction(
			"math:tan",
			ir.NewCode(c.Namespace, c.storage).
				DoubleTan("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.registerIRFunction(
			"math:acos",
			ir.NewCode(c.Namespace, c.storage).
				Copy(
					fmt.Sprintf("%s.acos.x", ArgPath),
					fmt.Sprintf("%s.x", VarPath),
				).
				DoubleAcos("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.registerIRFunction(
			"math:asin",
			ir.NewCode(c.Namespace, c.storage).
				Copy(
					fmt.Sprintf("%s.asin.x", ArgPath),
					fmt.Sprintf("%s.x", VarPath),
				).
				DoubleAsin("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.registerIRFunction(
			"math:atan",
			ir.NewCode(c.Namespace, c.storage).
				Copy(
					fmt.Sprintf("%s.atan.x", ArgPath),
					fmt.Sprintf("%s.x", VarPath),
				).
				DoubleAtan("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),

		// Rounding functions
		c.registerIRFunction(
			"math:floor",
			ir.NewCode(c.Namespace, c.storage).
				DoubleFloor("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.registerIRFunction(
			"math:ceil",
			ir.NewCode(c.Namespace, c.storage).
				DoubleCeil("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
		c.registerIRFunction(
			"math:round",
			ir.NewCode(c.Namespace, c.storage).
				DoubleRound("x", RET).
				Ret(),
			[]interfaces.TypedIdentifier{
				{Name: "x", Type: types.DoubleType},
			},
			types.DoubleType,
		),
	}
}

func (c *Compiler) baseFunctions() []interfaces.Function {

	funcs := make([]interfaces.Function, 0)

	initSource := ir.NewCode(c.Namespace, c.storage)
	if c.Config.CleanBeforeInit {
		initSource.Call(fmt.Sprintf("mcb:%s", path.Join(paths.Internal, "clean")))
	}
	initSource.Raw(fmt.Sprintf("scoreboard objectives add %s dummy", c.Namespace))
	if c.Config.Debug {
		initSource.Raw(fmt.Sprintf("scoreboard objectives setdisplay sidebar %s", c.Namespace))
	}
	initSource.Set(fmt.Sprintf("%s.%s", VarPath, CALL), nbt.NewInt(0))
	initSource.XLoad(CALL, CALL)
	initSource.SetArg("mcb:log", "text", nbt.NewString("MCB pack loaded"))
	initSource.Call("mcb:log")
	initSource.Call("load")
	initSource.Ret()
	funcs = append(
		funcs,
		c.registerIRFunction(
			fmt.Sprintf("mcb:%s", path.Join(paths.Internal, "init")),
			initSource,
			[]interfaces.TypedIdentifier{},
			types.VoidType,
		),
		c.registerIRFunction(
			fmt.Sprintf("mcb:%s", path.Join(paths.Internal, "clean")),
			ir.NewCode(c.Namespace, c.storage).
				Raw(fmt.Sprintf("scoreboard objectives remove %s", c.Namespace)).
				Remove(VarPath).
				Remove(StructPath).
				Remove(ArgPath).
				Ret(),
			[]interfaces.TypedIdentifier{},
			types.VoidType,
		),
	)
	return funcs
}
