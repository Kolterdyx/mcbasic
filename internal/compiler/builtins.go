package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/ir"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	"path"
)

func (c *Compiler) compileBuiltins() []interfaces.Function {
	return c.baseFunctions()
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
	initSource.Call(utils.FileSpecifier(c.Config.Project.Entrypoint, "load"))
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
