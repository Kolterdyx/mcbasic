package ops

import "fmt"

func (o *Op) arithmeticOperation(a string, b string, to string, operator string) string {
	cmd := ""
	cmd += o.MoveScore(a, RX)
	cmd += o.MoveScore(b, RB)
	cmd += fmt.Sprintf("scoreboard players operation %s %s %s %s %s\n", o.Namespace, RX, operator, o.Namespace, RB)
	cmd += o.MoveScore(RX, to)
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
