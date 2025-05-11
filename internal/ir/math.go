package ir

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

func (c *Code) intOperation(x, y, to, operator string) interfaces.IRCode {
	c.Load(x, RA)
	c.Load(y, RX)
	c.MathOp(operator)
	c.Store(RX, c.varPath(to))
	return c
}

func (c *Code) callGM1(op, x, to string) interfaces.IRCode {
	c.CopyArg(x, op, "x")
	c.Raw(fmt.Sprintf("function gm:%s with storage %s %s.%s", op, c.storage, ArgPath, op))
	c.XCopy("gm:io", "out", c.storage, c.varPath(to))
	return c
}

func (c *Code) callGM2(op, x, y, to string) interfaces.IRCode {
	c.CopyArg(x, op, "x")
	c.CopyArg(y, op, "y")
	c.Raw(fmt.Sprintf("function gm:%s with storage %s %s.%s", op, c.storage, ArgPath, op))
	c.XCopy("gm:io", "out", c.storage, c.varPath(to))
	return c
}

func (c *Code) IntAdd(x, y, to string) interfaces.IRCode {
	return c.intOperation(x, y, to, "+=")
}

func (c *Code) IntSub(x, y, to string) interfaces.IRCode {
	return c.intOperation(x, y, to, "-=")
}

func (c *Code) IntMul(x, y, to string) interfaces.IRCode {
	return c.intOperation(x, y, to, "*=")
}

func (c *Code) IntDiv(x, y, to string) interfaces.IRCode {
	return c.intOperation(x, y, to, "/=")
}

func (c *Code) IntMod(x, y, to string) interfaces.IRCode {
	return c.intOperation(x, y, to, "%=")
}

func (c *Code) DoubleAdd(x, y, to string) interfaces.IRCode {
	return c.callGM2("add", x, y, to)
}

func (c *Code) DoubleSub(x, y, to string) interfaces.IRCode {
	return c.callGM2("subtract", x, y, to)
}

func (c *Code) DoubleMul(x, y, to string) interfaces.IRCode {
	return c.callGM2("multiply", x, y, to)
}

func (c *Code) DoubleDiv(x, y, to string) interfaces.IRCode {
	return c.callGM2("divide", x, y, to)
}

func (c *Code) DoubleMod(x, y, to string) interfaces.IRCode {
	return c.callGM2("modulo", x, y, to)
}

func (c *Code) DoubleSqrt(x, to string) interfaces.IRCode {
	return c.callGM1("sqrt", x, to)
}

func (c *Code) DoubleCos(x, to string) interfaces.IRCode {
	return c.callGM1("cos", x, to)
}

func (c *Code) DoubleSin(x, to string) interfaces.IRCode {
	return c.callGM1("sin", x, to)
}

func (c *Code) DoubleTan(x, to string) interfaces.IRCode {
	return c.callGM1("tan", x, to)
}

func (c *Code) DoubleAcos(x, to string) interfaces.IRCode {
	return c.callGM1("arccos", x, to)
}

func (c *Code) DoubleAsin(x, to string) interfaces.IRCode {
	return c.callGM1("arcsin", x, to)
}

func (c *Code) DoubleAtan(x, to string) interfaces.IRCode {
	return c.callGM1("arctan", x, to)
}

func (c *Code) DoubleAtan2(x, y, to string) interfaces.IRCode {
	return c.callGM2("arctan2", x, y, to)
}

func (c *Code) DoubleFloor(x, to string) interfaces.IRCode {
	return c.callGM1("floor", x, to)
}

func (c *Code) DoubleCeil(x, to string) interfaces.IRCode {
	return c.callGM1("ceil", x, to)
}

func (c *Code) DoubleRound(x, to string) interfaces.IRCode {
	return c.callGM1("round", x, to)
}
