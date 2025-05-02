package il

func (c *Compiler) intOperation(x, y, to, operator string) (cmd string) {
	cmd += c.Load(c.varPath(x), RA)
	cmd += c.Load(c.varPath(y), RB)
	cmd += c.MathOp(operator)
	cmd += c.Store(RX, c.varPath(to))
	return
}

func (c *Compiler) IntAdd(a, b, to string) string {
	return c.intOperation(a, b, to, "add")
}

func (c *Compiler) IntSub(a, b, to string) string {
	return c.intOperation(a, b, to, "sub")
}

func (c *Compiler) IntMul(a, b, to string) string {
	return c.intOperation(a, b, to, "mul")
}

func (c *Compiler) IntDiv(a, b, to string) string {
	return c.intOperation(a, b, to, "div")
}

func (c *Compiler) IntMod(a, b, to string) string {
	return c.intOperation(a, b, to, "mod")
}

func (c *Compiler) DoubleAdd(a, b, to string) string {
	return ""
}

func (c *Compiler) DoubleSub(a, b, to string) string {
	return ""
}

func (c *Compiler) DoubleMul(a, b, to string) string {
	return ""
}

func (c *Compiler) DoubleDiv(a, b, to string) string {
	return ""
}

func (c *Compiler) DoubleMod(a, b, to string) string {
	return ""
}
