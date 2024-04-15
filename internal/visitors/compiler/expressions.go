package compiler

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
)

func (c *Compiler) VisitLiteral(expr expressions.LiteralExpr) interface{} {
	return c.opHandler.MoveConst(expr.Value.(string), ops.RX)
}

func (c *Compiler) VisitBinary(expr expressions.BinaryExpr) interface{} {
	cmd := ""

	cmd += expr.Left.Accept(c).(string)
	c.opHandler.Move(ops.RX, ops.RA)
	cmd += expr.Right.Accept(c).(string)
	c.opHandler.Move(ops.RX, ops.RB)
	if expr.ReturnType() == expressions.NumberType {
		switch expr.Operator.Type {
		case tokens.Plus:
			return c.opHandler.Add(ops.RA, ops.RB, ops.RX)
		case tokens.Minus:
			return c.opHandler.Sub(ops.RA, ops.RB, ops.RX)
		case tokens.Star:
			return c.opHandler.Mul(ops.RA, ops.RB, ops.RX)
		case tokens.Slash:
			return c.opHandler.Div(ops.RA, ops.RB, ops.RX)
		case tokens.Percent:
			return c.opHandler.Mod(ops.RA, ops.RB, ops.RX)
		default:
			panic("Unknown operator")
		}
	} else if expr.ReturnType() == expressions.StringType {
		if expr.Operator.Type == tokens.Plus {
			return c.opHandler.Concat(ops.RA, ops.RB, ops.RX)
		} else {
			panic("Unknown operator")
		}
	} else {
		panic("Unknown return type")
	}
	return cmd
}

func (c *Compiler) VisitVariable(expr expressions.VariableExpr) interface{} {
	return c.opHandler.Move(expr.Name.Lexeme, ops.RX)
}
