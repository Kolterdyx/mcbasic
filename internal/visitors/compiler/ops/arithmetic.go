package ops

import (
	"fmt"
	"math"
	"strconv"
)

func (o *Op) arithmeticOperation(a string, b string, to string, operator string) string {
	cmd := ""
	cmd += o.MoveScore(Cs(a), Cs(RX))
	cmd += o.MoveScore(Cs(b), Cs(RB))
	cmd += fmt.Sprintf("scoreboard players operation %s %s %s %s %s\n", Cs(RX), o.Namespace, operator, Cs(RB), o.Namespace)
	cmd += o.LoadScore(Cs(RX), Cs(to))
	return cmd
}

func (o *Op) Add(a string, b string, to string) string {
	return o.arithmeticOperation(a, b, to, "+=")
}

func (o *Op) Sub(a string, b string, to string) string {
	return o.arithmeticOperation(a, b, to, "-=")
}

func (o *Op) Mul(a string, b string, to string) string {
	return o.arithmeticOperation(a, b, to, "*=")
}

func (o *Op) Div(a string, b string, to string) string {
	return o.arithmeticOperation(a, b, to, "/=")
}

func (o *Op) Mod(a string, b string, to string) string {
	return o.arithmeticOperation(a, b, to, "%=")
}

func (o *Op) Scale(value string, scale string, to string) string {
	return fmt.Sprintf("execute store result storage %s:%s %s double %s run data get storage %s:%s %s\n", o.Namespace, VarPath, Cs(to), scale, o.Namespace, VarPath, Cs(value))
}

func (o *Op) FixedAdd(a string, b string, to string) string {
	return o.Add(Cs(a), Cs(b), Cs(to))
}

func (o *Op) FixedSub(a string, b string, to string) string {
	return o.Sub(Cs(a), Cs(b), Cs(to))
}

func (o *Op) FixedMul(a string, b string, to string) string {
	cmd := ""
	cmd += o.Mul(Cs(a), Cs(b), Cs(to))
	invn := strconv.FormatFloat(1/math.Pow(10, float64(o.FixedPointPrecision)), 'f', -1, 64)
	cmd += o.Scale(Cs(to), invn, Cs(to))
	return cmd
}

func (o *Op) FixedDiv(a string, b string, to string) string {
	cmd := ""
	n := strconv.FormatFloat(math.Pow(10, float64(o.FixedPointPrecision)), 'f', -1, 64)
	// To avoid weird number boundary shenanigans, we can just multiply the numerator by the inverse of the denominator,
	// instead of scaling up the numerator
	// a / b = a * (1/b)

	cmd += o.MoveConst("1.0", Cs(RB))
	cmd += o.Scale(Cs(RB), n, Cs(RB))
	cmd += o.Div(Cs(RB), Cs(b), Cs(RB))
	cmd += o.FixedMul(a, Cs(RB), to)

	return cmd
}
