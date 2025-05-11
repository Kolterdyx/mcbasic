package ir

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	log "github.com/sirupsen/logrus"
)

func (c *Code) intOperation(x, y, to, operator string) interfaces.IRCode {
	c.Load(x, RA)
	c.Load(y, RX)
	c.MathOp(operator)
	c.Store(RX, c.varPath(to))
	return c
}

func (c *Code) callGM1(x, to string) interfaces.IRCode {
	panic("not implemented")
}

func (c *Code) callGM2(x, y, to string) interfaces.IRCode {
	panic("not implemented")
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
	log.Warnln("DoubleAdd has not been implemented yet")
	return c
}

func (c *Code) DoubleSub(x, y, to string) interfaces.IRCode {
	log.Warnln("DoubleSub has not been implemented yet")
	return c
}

func (c *Code) DoubleMul(x, y, to string) interfaces.IRCode {
	log.Warnln("DoubleMul has not been implemented yet")
	return c
}

func (c *Code) DoubleDiv(x, y, to string) interfaces.IRCode {
	log.Warnln("DoubleDiv has not been implemented yet")
	return c
}

func (c *Code) DoubleMod(x, y, to string) interfaces.IRCode {
	log.Warnln("DoubleMod has not been implemented yet")
	return c
}

func (c *Code) DoubleSqrt(x, to string) interfaces.IRCode {
	log.Warnln("DoubleSqrt has not been implemented yet")
	return c
}

func (c *Code) DoubleCos(x, to string) interfaces.IRCode {
	log.Warnln("DoubleCos has not been implemented yet")
	return c
}

func (c *Code) DoubleSin(x, to string) interfaces.IRCode {
	log.Warnln("DoubleSin has not been implemented yet")
	return c
}

func (c *Code) DoubleTan(x, to string) interfaces.IRCode {
	log.Warnln("DoubleTan has not been implemented yet")
	return c
}

func (c *Code) DoubleAcos(x, to string) interfaces.IRCode {
	log.Warnln("DoubleAcos has not been implemented yet")
	return c
}

func (c *Code) DoubleAsin(x, to string) interfaces.IRCode {
	log.Warnln("DoubleAsin has not been implemented yet")
	return c
}

func (c *Code) DoubleAtan(x, to string) interfaces.IRCode {
	log.Warnln("DoubleAtan has not been implemented yet")
	return c
}

func (c *Code) DoubleAtan2(x, y, to string) interfaces.IRCode {
	log.Warnln("DoubleAtan2 has not been implemented yet")
	return c
}

func (c *Code) DoubleFloor(x, to string) interfaces.IRCode {
	log.Warnln("DoubleFloor has not been implemented yet")
	return c
}

func (c *Code) DoubleCeil(x, to string) interfaces.IRCode {
	log.Warnln("DoubleCeil has not been implemented yet")
	return c
}

func (c *Code) DoubleRound(x, to string) interfaces.IRCode {
	log.Warnln("DoubleRound has not been implemented yet")
	return c
}
