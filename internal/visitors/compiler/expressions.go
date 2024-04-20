package compiler

import (
	"fmt"
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
		if expr.Operator.Type == tokens.Plus {
			cmd += c.opHandler.Concat(ops.Cs(regRa), ops.Cs(regRb), ops.Cs(ops.RX))
		} else {
			c.error(expr.SourceLocation, "Invalid operator for strings")
		}
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

func (c *Compiler) VisitGrouping(expr expressions.GroupingExpr) interface{} {
	return expr.Expression.Accept(c).(string)
}

func (c *Compiler) VisitLogical(expr expressions.LogicalExpr) interface{} {
	leftSide := ""
	rightSide := ""

	regRa := c.newRegister(ops.RA)
	regRb := c.newRegister(ops.RB)

	cmd := ""

	leftSide += "### Logical operation left side ###\n"
	leftSide += expr.Left.Accept(c).(string)
	leftSide += c.opHandler.Move(ops.Cs(ops.RX), regRa)
	rightSide += "### Logical operation right side ###\n"
	rightSide += expr.Right.Accept(c).(string)
	rightSide += c.opHandler.Move(ops.Cs(ops.RX), regRb)
	rightSide += c.opHandler.MoveScore(regRb, regRb)

	cmd += leftSide
	switch expr.Operator.Type {
	case tokens.And:
		// If left side is false, return false
		evalRightSide := ""
		cmd += c.opHandler.MoveScore(regRa, regRa)
		cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), true, c.opHandler.MoveConst("0", ops.Cs(ops.RX)))
		evalRightSide += rightSide
		evalRightSide += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), true, c.opHandler.MoveConst("0", ops.Cs(ops.RX)))
		evalRightSide += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), false, c.opHandler.Move(regRb, ops.Cs(ops.RX)))
		cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), false, evalRightSide)
	case tokens.Or:
		// If left side is true, return true
		evalRightSide := ""
		cmd += c.opHandler.MoveScore(regRa, regRa)
		cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), false, c.opHandler.Move(regRa, ops.Cs(ops.RX)))
		evalRightSide += rightSide
		evalRightSide += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), true, c.opHandler.MoveConst("0", ops.Cs(ops.RX)))
		evalRightSide += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), false, c.opHandler.Move(regRb, ops.Cs(ops.RX)))
		cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), true, evalRightSide)
	default:
		c.error(expr.SourceLocation, "Invalid operator for logical expressions")
	}

	return cmd
}

func (c *Compiler) VisitSlice(expr expressions.SliceExpr) interface{} {
	cmd := ""
	cmd += "### Slice operation ###\n"
	cmd += expr.StartIndex.Accept(c).(string)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(ops.RA))
	if expr.EndIndex == nil {
		cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(ops.RB))
		cmd += c.opHandler.Inc(ops.Cs(ops.RB))
	} else {
		cmd += expr.EndIndex.Accept(c).(string)
		cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(ops.RB))
	}
	cmd += expr.TargetExpr.Accept(c).(string)
	cmd += c.opHandler.TraceStorage("mcb:vars", ops.Cs(ops.RX), "RX")
	cmd += c.opHandler.Slice(ops.Cs(ops.RX), ops.Cs(ops.RA), ops.Cs(ops.RB), ops.Cs(ops.RX))
	return cmd
}
