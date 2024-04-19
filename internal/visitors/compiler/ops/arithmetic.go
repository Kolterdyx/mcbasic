package ops

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal"
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

func (o *Op) FixedAdd(a string, b string, to string) string {
	cmd := ""
	aw := Cs(a) + ".whole"
	af := Cs(a) + ".fract"
	bw := Cs(b) + ".whole"
	bf := Cs(b) + ".fract"
	tw := Cs(to) + ".whole"
	tf := Cs(to) + ".fract"
	cmd += o.Move(aw, Cs(RX))
	cmd += o.Move(bw, Cs(RB))
	cmd += o.Add(RX, RB, tw)
	cmd += o.Move(af, Cs(RX))
	cmd += o.Move(bf, Cs(RB))
	cmd += o.Add(RX, RB, tf)
	// Carry
	precisionMagnitude := int(math.Pow(10, internal.FixedPointMagnitude))
	tmpReg := Cs(RX) + "tmp"
	cmd += o.MoveScore(tf, tmpReg)
	carry := ""
	carry += o.MoveConst("1", Cs(RX))
	carry += o.Add(tw, Cs(RX), tw)
	carry += o.MoveConst(strconv.Itoa(precisionMagnitude), Cs(RX))
	carry += o.Sub(tf, Cs(RX), tf)
	cmd += o.ExecCond(fmt.Sprintf("score %s %s matches %d..", tmpReg, o.Namespace, precisionMagnitude), true, carry)
	return cmd
}

func (o *Op) FixedSub(a string, b string, to string) string {
	cmd := ""
	aw := Cs(a) + ".whole"
	af := Cs(a) + ".fract"
	bw := Cs(b) + ".whole"
	bf := Cs(b) + ".fract"
	tw := Cs(to) + ".whole"
	tf := Cs(to) + ".fract"
	cmd += o.MoveFixedConst("0.0", Cs(to))
	cmd += o.Move(aw, Cs(RX))
	cmd += o.Move(bw, Cs(RB))
	cmd += o.Sub(RX, RB, tw)
	cmd += o.Move(af, Cs(RX))
	cmd += o.Move(bf, Cs(RB))
	cmd += o.Sub(RX, RB, tf)
	// Borrow
	precisionMagnitude := int(math.Pow(10, internal.FixedPointMagnitude))
	tmpReg := Cs(RX) + "tmp"
	cmd += o.MoveScore(tf, tmpReg)
	borrow := ""
	borrow += o.MoveConst("1", Cs(RX))
	borrow += o.Sub(tw, Cs(RX), tw)
	borrow += o.MoveConst(strconv.Itoa(precisionMagnitude), Cs(RX))
	borrow += o.Add(tf, Cs(RX), tf)
	cmd += o.ExecCond(fmt.Sprintf("score %s %s matches ..-1", tmpReg, o.Namespace), true, borrow)
	return cmd
}
