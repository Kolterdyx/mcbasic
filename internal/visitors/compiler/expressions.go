package compiler

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	"strconv"
)

func (c *Compiler) VisitLiteral(expr expressions.LiteralExpr) interface{} {
	switch expr.ReturnType() {
	case expressions.IntType:
		return c.opHandler.MoveConst(expr.Value.(string), ops.Cs(ops.RX))
	case expressions.StringType:
		return c.opHandler.MoveConst(strconv.Quote(expr.Value.(string)), ops.Cs(ops.RX))
	case expressions.FixedType:
		return c.opHandler.MoveConst(expr.Value.(string), ops.Cs(ops.RX))
	default:
		c.error(expr.SourceLocation, "Invalid type in literal expression")
	}
	return ""
}

func (c *Compiler) VisitBinary(expr expressions.BinaryExpr) interface{} {
	cmd := ""

	regRa := c.newRegister(ops.RA)
	regRb := c.newRegister(ops.RB)

	cmd += "### Binary operation left side ###\n"
	cmd += expr.Left.Accept(c).(string)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regRa))
	cmd += "### Binary operation right side ###\n"
	cmd += expr.Right.Accept(c).(string)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regRb))

	cmd += "### Binary operation ###\n"

	switch expr.Operator.Type {
	case tokens.EqualEqual, tokens.BangEqual, tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual:
		cmd += c.Compare(expr, ops.Cs(regRa), ops.Cs(regRb), ops.Cs(ops.RX))
		return cmd
	default:
		// Do nothing
	}

	switch expr.ReturnType() {
	case expressions.IntType:
		switch expr.Operator.Type {
		case tokens.Plus:
			cmd += c.opHandler.Add(regRa, regRb, ops.RX)
		case tokens.Minus:
			cmd += c.opHandler.Sub(regRa, regRb, ops.RX)
		case tokens.Star:
			cmd += c.opHandler.Mul(regRa, regRb, ops.RX)
		case tokens.Slash:
			cmd += c.opHandler.Div(regRa, regRb, ops.RX)
		case tokens.Percent:
			cmd += c.opHandler.Mod(regRa, regRb, ops.RX)
		default:
			c.error(expr.SourceLocation, "Invalid operator for integers")
		}
	case expressions.FixedType:
		switch expr.Operator.Type {
		case tokens.Plus:
			cmd += c.opHandler.FixedAdd(regRa, regRb, ops.RX)
		case tokens.Minus:
			cmd += c.opHandler.FixedSub(regRa, regRb, ops.RX)
		case tokens.Star:
			cmd += c.opHandler.FixedMul(regRa, regRb, ops.RX)
		case tokens.Slash:
			cmd += c.opHandler.FixedDiv(regRa, regRb, ops.RX)
		default:
			c.error(expr.SourceLocation, "Invalid operator for fixed point numbers")
		}

	case expressions.StringType:
		c.error(expr.SourceLocation, "String operations are not supported yet")
		//if expr.Operator.Type == tokens.Plus {
		//	return c.opHandler.Concat(ops.RA, ops.RB, ops.RX)
		//} else {
		//	panic("Unknown operator")
		//}
	default:
		c.error(expr.SourceLocation, "Invalid type combination in binary operation")
	}
	return cmd
}

func (c *Compiler) VisitVariable(expr expressions.VariableExpr) interface{} {
	return c.opHandler.Move(ops.Cs(expr.Name.Lexeme), ops.Cs(ops.RX))
}

func (c *Compiler) VisitFunctionCall(expr expressions.FunctionCallExpr) interface{} {
	cmd := ""
	for i, arg := range expr.Arguments {
		cmd += arg.Accept(c).(string)
		argName := c.functions[expr.Name.Lexeme].Args[i].Name
		cmd += c.opHandler.LoadArg(expr.Name.Lexeme, argName, ops.Cs(ops.RX))
	}
	if expr.ReturnType() != expressions.VoidType {
		cmd += c.opHandler.Call(expr.Name.Lexeme, ops.RX)
	} else {
		cmd += c.opHandler.Call(expr.Name.Lexeme, "")
	}
	return cmd
}

func (c *Compiler) VisitUnary(expr expressions.UnaryExpr) interface{} {
	cmd := ""
	switch expr.ReturnType() {
	case expressions.IntType:
		switch expr.Operator.Type {
		case tokens.Minus:
			cmd += expressions.BinaryExpr{
				Left: expressions.LiteralExpr{Value: "0", SourceLocation: expr.SourceLocation, ValueType: expressions.IntType},
				Operator: tokens.Token{
					Type: tokens.Minus,
				},
				Right: expr.Expression,
			}.Accept(c).(string)
		case tokens.Bang:
			cmd += expr.Expression.Accept(c).(string)
			cmd += c.opHandler.NegateNumber(ops.Cs(ops.RX))
		default:
			c.error(expr.SourceLocation, "Invalid operator for integers")
		}
	default:
		c.error(expr.SourceLocation, "Invalid type in unary expression")
	}
	return cmd
}
