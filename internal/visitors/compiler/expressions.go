package compiler

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	log "github.com/sirupsen/logrus"
	"slices"
)

func (c *Compiler) VisitLiteral(expr expressions.LiteralExpr) interface{} {
	if expr.ValueType == expressions.NumberType {
		return c.opHandler.RegLoad(expr.Value.(string), ops.RX)
	}
	return c.opHandler.Set(ops.RX, expr.Value.(string))
}

func (c *Compiler) VisitVariable(expr expressions.VariableExpr) interface{} {
	return c.opHandler.RegWrite(expr.Name.Lexeme, ops.RX)
}

func (c *Compiler) VisitUnary(expr expressions.UnaryExpr) interface{} {
	cmd := expr.Expression.Accept(c).(string)
	cmd += c.opHandler.RegShift(ops.RX, ops.RA)
	switch expr.Operator.Type {
	case tokens.Minus:
		cmd += c.opHandler.RegLoad("0", ops.RX)
		cmd += c.opHandler.Sub(ops.RX, ops.RA)
	case tokens.Bang:
		cmd += c.opHandler.RegLoad("1", ops.RX)
		cmd += c.opHandler.Sub(ops.RX, ops.RA)
	default:
		log.Fatalln("Unknown unary operator")
	}
	return cmd
}

func (c *Compiler) VisitBinary(expr expressions.BinaryExpr) interface{} {
	arithmeticOperators := []tokens.TokenType{
		tokens.Plus,
		tokens.Minus,
		tokens.Star,
		tokens.Slash,
		tokens.Percent,
	}
	if slices.Contains(arithmeticOperators, expr.Operator.Type) {
		if expr.Left.Type() == expressions.LiteralExprType {
			left := expr.Left.(expressions.LiteralExpr)
			if left.ValueType != expressions.NumberType {
				c.error(left.SourceLocation, "Left side of binary arithmetic expression is not a number")
			}
		}
		if expr.Right.Type() == expressions.LiteralExprType {
			right := expr.Right.(expressions.LiteralExpr)
			if right.ValueType != expressions.NumberType {
				c.error(right.SourceLocation, "Right side of binary arithmetic expression is not a number")
			}
		}
	}

	cmd := expr.Left.Accept(c).(string)
	cmd += c.opHandler.RegShift(ops.RX, ops.RA)
	cmd += expr.Right.Accept(c).(string)
	cmd += c.opHandler.RegShift(ops.RX, ops.RB)
	switch expr.Operator.Type {
	case tokens.Plus:
		cmd += c.opHandler.Add(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.Minus:
		cmd += c.opHandler.Sub(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.Star:
		cmd += c.opHandler.Mul(ops.RA, ops.RB)
		cmd += c.opHandler.Div(ops.RA, ops.RCF) // Compensate for fixed-point shift
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.Slash:
		cmd += c.opHandler.Div(ops.RB, ops.RCF) // Compensate for fixed-point shift
		cmd += c.opHandler.Div(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.Percent:
		cmd += c.opHandler.Mod(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	default:
		log.Debug("Not an arithmetic operator")
	}

	switch expr.Operator.Type {
	case tokens.EqualEqual:
		cmd += c.opHandler.Eq(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.BangEqual:
		cmd += c.opHandler.Neq(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.Less:
		cmd += c.opHandler.Lt(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.LessEqual:
		cmd += c.opHandler.Lte(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.Greater:
		cmd += c.opHandler.Gt(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.GreaterEqual:
		cmd += c.opHandler.Gte(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	default:
		log.Debug("Not a comparison operator")
	}

	switch expr.Operator.Type {
	case tokens.And:
		cmd += c.opHandler.And(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	case tokens.Or:
		cmd += c.opHandler.Or(ops.RA, ops.RB)
		cmd += c.opHandler.RegShift(ops.RA, ops.RX)
		return cmd
	default:
		log.Debug("Not a logical operator")
	}

	log.Fatalln("Unknown binary operator")
	return ""
}

func (c *Compiler) VisitGrouping(expr expressions.GroupingExpr) interface{} {
	return expr.Expression.Accept(c)
}

func (c *Compiler) VisitFunctionCall(expr expressions.FunctionCallExpr) interface{} {
	cmd := ""
	for i, arg := range expr.Arguments {
		cmd += arg.Accept(c).(string)
		cmd += c.opHandler.RegSave(ops.RX, ops.RX)
		cmd += c.opHandler.ArgLoad(expr.Name.Lexeme, c.functionArgs[expr.Name.Lexeme][i], ops.RX)
	}
	cmd += c.opHandler.Call(expr.Name.Lexeme)
	return cmd
}