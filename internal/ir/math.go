package ir

import log "github.com/sirupsen/logrus"

func (c *Compiler) intOperation(x, y, to, operator string) (cmd string) {
	cmd += c.Load(x, RA)
	cmd += c.Load(y, RX)
	cmd += c.MathOp(operator)
	cmd += c.Store(RX, c.varPath(to))
	return
}

func (c *Compiler) callGM1(x, to string) string {
	panic("not implemented")
}

func (c *Compiler) callGM2(x, y, to string) string {
	panic("not implemented")
}

func (c *Compiler) IntAdd(x, y, to string) string {
	return c.intOperation(x, y, to, "+=")
}

func (c *Compiler) IntSub(x, y, to string) string {
	return c.intOperation(x, y, to, "-=")
}

func (c *Compiler) IntMul(x, y, to string) string {
	return c.intOperation(x, y, to, "*=")
}

func (c *Compiler) IntDiv(x, y, to string) string {
	return c.intOperation(x, y, to, "/=")
}

func (c *Compiler) IntMod(x, y, to string) string {
	return c.intOperation(x, y, to, "%=")
}

func (c *Compiler) DoubleAdd(x, y, to string) string {
	log.Warnln("DoubleAdd has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleSub(x, y, to string) string {
	log.Warnln("DoubleSub has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleMul(x, y, to string) string {
	log.Warnln("DoubleMul has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleDiv(x, y, to string) string {
	log.Warnln("DoubleDiv has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleMod(x, y, to string) string {
	log.Warnln("DoubleMod has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleSqrt(x, to string) string {
	log.Warnln("DoubleSqrt has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleCos(x, to string) string {
	log.Warnln("DoubleCos has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleSin(x, to string) string {
	log.Warnln("DoubleSin has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleTan(x, to string) string {
	log.Warnln("DoubleTan has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleAcos(x, to string) string {
	log.Warnln("DoubleAcos has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleAsin(x, to string) string {
	log.Warnln("DoubleAsin has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleAtan(x, to string) string {
	log.Warnln("DoubleAtan has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleAtan2(x, y, to string) string {
	log.Warnln("DoubleAtan2 has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleFloor(x, to string) string {
	log.Warnln("DoubleFloor has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleCeil(x, to string) string {
	log.Warnln("DoubleCeil has not been implemented yet")
	return ""
}

func (c *Compiler) DoubleRound(x, to string) string {
	log.Warnln("DoubleRound has not been implemented yet")
	return ""
}
