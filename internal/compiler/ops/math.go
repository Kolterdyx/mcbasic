package ops

import (
	"fmt"
)

func (o *Op) integerOperation(x, y, to, operator string) string {
	cmd := ""
	cmd += o.MoveScore(Cs(x), Cs(RX))
	cmd += o.MoveScore(Cs(y), Cs(RB))
	cmd += fmt.Sprintf("scoreboard players operation %s %s %s %s %s\n", Cs(RX), o.Namespace, operator, Cs(RB), o.Namespace)
	cmd += o.LoadScore(Cs(RX), Cs(to))
	return cmd
}

func (o *Op) MoveGMResult(to string) string {
	return o.MoveRaw("gm:io", "out", fmt.Sprintf("%s:data", o.Namespace), fmt.Sprintf("%s.%s", VarPath, to))
}

func (o *Op) CallGM1(op, x, to string) string {
	return o.CallFunction(op, map[string]string{
		"x": x,
	}, "") + o.MoveGMResult(to)
}

func (o *Op) CallGM2(op, x, y, to string) string {
	return o.CallFunction(op, map[string]string{
		"x": x,
		"y": y,
	}, "") + o.MoveGMResult(to)
}

func (o *Op) Add(x, y, to string) string {
	return o.integerOperation(x, y, to, "+=")
}

func (o *Op) Sub(x, y, to string) string {
	return o.integerOperation(x, y, to, "-=")
}

func (o *Op) Mul(x, y, to string) string {
	return o.integerOperation(x, y, to, "*=")
}

func (o *Op) Div(x, y, to string) string {
	return o.integerOperation(x, y, to, "/=")
}

func (o *Op) Mod(x, y, to string) string {
	return o.integerOperation(x, y, to, "%=")
}

func (o *Op) Scale(value, scale, to string) string {
	return fmt.Sprintf("execute store result storage %s:%s %s double %s run data get storage %s:%s %s\n", o.Namespace, VarPath, Cs(to), scale, o.Namespace, VarPath, Cs(value))
}

func (o *Op) DoubleAdd(x, y, to string) string {
	return o.CallGM2("gm:add", x, y, Cs(to))
}

func (o *Op) DoubleSub(x, y, to string) string {
	return o.CallGM2("gm:subtract", x, y, Cs(to))
}

func (o *Op) DoubleMul(x, y, to string) string {
	return o.CallGM2("gm:multiply", x, y, Cs(to))
}

func (o *Op) DoubleDiv(x, y, to string) string {
	return o.CallGM2("gm:divide", x, y, Cs(to))
}

func (o *Op) DoubleMod(x, y, to string) string {
	return o.CallGM2("gm:modulo", x, y, Cs(to))
}

func (o *Op) DoubleSqrt(x, to string) string {
	return o.CallGM1("gm:sqrt", x, to)
}

func (o *Op) DoubleSin(x, to string) string {
	return o.CallGM1("gm:sin", x, to)
}

func (o *Op) DoubleCos(x, to string) string {
	return o.CallGM1("gm:cos", x, to)
}

func (o *Op) DoubleTan(x, to string) string {
	return o.CallGM1("gm:tan", x, to)
}

func (o *Op) DoubleAsin(x, to string) string {
	return o.CallGM1("gm:arcsin", x, to)
}

func (o *Op) DoubleAcos(x, to string) string {
	return o.CallGM1("gm:arccos", x, to)
}

func (o *Op) DoubleAtan(x, to string) string {
	return o.CallGM1("gm:arctan", x, to)
}

func (o *Op) DoubleRound(x, to string) string {
	return o.CallGM1("gm:round", x, to)
}

func (o *Op) DoubleFloor(x, to string) string {
	return o.CallGM1("gm:floor", x, to)
}

func (o *Op) DoubleCeil(x, to string) string {
	return o.CallGM1("gm:ceil", x, to)
}
