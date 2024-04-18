package ops

import "fmt"

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
	af := Cs(a) + ".fixed"
	bw := Cs(b) + ".whole"
	bf := Cs(b) + ".fixed"
	tw := Cs(to) + ".whole"
	tf := Cs(to) + ".fixed"
	cmd += o.Move(aw, Cs(RX))
	cmd += o.Move(bw, Cs(RB))
	cmd += o.Add(RX, RB, tw)
	cmd += o.Move(af, Cs(RX))
	cmd += o.Move(bf, Cs(RB))
	cmd += o.Add(RX, RB, tf)
	// Carry
	//cmd += o.MoveScore(Cs(tf), Cs(RB))
	//carry := ""
	//carry += o.MoveConst("1", Cs(RX))
	//carry += o.Add(RX, RB, to)
	//cmd += o.ExecCond(fmt.Sprintf("score %s %s matches 10..", Cs(RB), o.Namespace), true, carry)
	return cmd
}
