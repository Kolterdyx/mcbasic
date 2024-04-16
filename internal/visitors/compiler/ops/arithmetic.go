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
