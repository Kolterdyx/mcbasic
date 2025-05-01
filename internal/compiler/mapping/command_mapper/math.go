package command_mapper

import (
	"fmt"
)

func (c *CommandMapper) integerOperation(x, y, to, operator string) string {
	cmd := ""
	cmd += c.MoveScore(c.Cs(x), c.Cs(c.RX))
	cmd += c.MoveScore(c.Cs(y), c.Cs(c.RB))
	cmd += fmt.Sprintf("scoreboard players operation %s %s %s %s %s\n", c.Cs(c.RX), c.Namespace, operator, c.Cs(c.RB), c.Namespace)
	cmd += c.LoadScore(c.Cs(c.RX), c.Cs(to))
	return cmd
}

func (c *CommandMapper) moveGMResult(to string) string {
	return c.MoveRaw("gm:io", "out", fmt.Sprintf("%s:data", c.Namespace), fmt.Sprintf("%s.%s", c.VarPath, to))
}

func (c *CommandMapper) callGM1(op, x, to string) string {
	cmd := c.LoadArgs(op, map[string]string{
		"x": x,
	})
	cmd += c.Call(op, "")
	cmd += c.moveGMResult(to)
	return cmd
}

func (c *CommandMapper) callGM2(op, x, y, to string) string {
	cmd := c.LoadArgs(op, map[string]string{
		"x": x,
		"y": y,
	})
	cmd += c.Call(op, "")
	cmd += c.moveGMResult(to)
	return cmd
}

func (c *CommandMapper) IntAdd(x, y, to string) string {
	return c.integerOperation(x, y, to, "+=")
}

func (c *CommandMapper) IntSub(x, y, to string) string {
	return c.integerOperation(x, y, to, "-=")
}

func (c *CommandMapper) IntMul(x, y, to string) string {
	return c.integerOperation(x, y, to, "*=")
}

func (c *CommandMapper) IntDiv(x, y, to string) string {
	return c.integerOperation(x, y, to, "/=")
}

func (c *CommandMapper) IntMod(x, y, to string) string {
	return c.integerOperation(x, y, to, "%=")
}

func (c *CommandMapper) DoubleAdd(x, y, to string) string {
	return c.callGM2("gm:add", x, y, c.Cs(to))
}

func (c *CommandMapper) DoubleSub(x, y, to string) string {
	return c.callGM2("gm:subtract", x, y, c.Cs(to))
}

func (c *CommandMapper) DoubleMul(x, y, to string) string {
	return c.callGM2("gm:multiply", x, y, c.Cs(to))
}

func (c *CommandMapper) DoubleDiv(x, y, to string) string {
	return c.callGM2("gm:divide", x, y, c.Cs(to))
}

func (c *CommandMapper) DoubleMod(x, y, to string) string {
	return c.callGM2("gm:modulo", x, y, c.Cs(to))
}

func (c *CommandMapper) DoubleSqrt(x, to string) string {
	return c.callGM1("gm:sqrt", x, to)
}

func (c *CommandMapper) DoubleSin(x, to string) string {
	return c.callGM1("gm:sin", x, to)
}

func (c *CommandMapper) DoubleCos(x, to string) string {
	return c.callGM1("gm:cos", x, to)
}

func (c *CommandMapper) DoubleTan(x, to string) string {
	return c.callGM1("gm:tan", x, to)
}

func (c *CommandMapper) DoubleAsin(x, to string) string {
	return c.callGM1("gm:arcsin", x, to)
}

func (c *CommandMapper) DoubleAcos(x, to string) string {
	return c.callGM1("gm:arccos", x, to)
}

func (c *CommandMapper) DoubleAtan(x, to string) string {
	return c.callGM1("gm:arctan", x, to)
}

func (c *CommandMapper) DoubleRound(x, to string) string {
	return c.callGM1("gm:round", x, to)
}

func (c *CommandMapper) DoubleFloor(x, to string) string {
	return c.callGM1("gm:floor", x, to)
}

func (c *CommandMapper) DoubleCeil(x, to string) string {
	return c.callGM1("gm:ceil", x, to)
}
