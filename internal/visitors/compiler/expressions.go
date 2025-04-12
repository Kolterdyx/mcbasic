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
	case expressions.DoubleType:
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
	case expressions.DoubleType:
		switch expr.Operator.Type {
		case tokens.Plus:
			cmd += c.opHandler.DoubleAdd(regRa, regRb, ops.RX)
		case tokens.Minus:
			cmd += c.opHandler.DoubleSub(regRa, regRb, ops.RX)
		case tokens.Star:
			cmd += c.opHandler.DoubleMul(regRa, regRb, ops.RX)
		case tokens.Slash:
			cmd += c.opHandler.DoubleDiv(regRa, regRb, ops.RX)
		default:
			c.error(expr.SourceLocation, "Invalid operator for double numbers")
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
	regIndexStart := c.newRegister(ops.RA)
	regIndexEnd := c.newRegister(ops.RB)

	cmd := ""
	cmd += "### SliceString operation ###\n"
	cmd += expr.StartIndex.Accept(c).(string)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regIndexStart))

	if expr.EndIndex == nil {
		cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regIndexEnd))
		cmd += c.opHandler.Inc(ops.Cs(regIndexEnd))
	} else {
		cmd += expr.EndIndex.Accept(c).(string)
		cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regIndexEnd))
	}
	cmd += expr.TargetExpr.Accept(c).(string)

	// Check index bounds
	lenReg := c.newRegister(ops.RX)

	cmd += c.opHandler.SizeString(ops.Cs(ops.RX), ops.Cs(lenReg))

	// If any of the indexes are negative, add the length of the string to them
	cmd += c.opHandler.ExecCond(
		fmt.Sprintf("score %s %s matches ..-1", ops.Cs(regIndexStart), c.Namespace),
		true,
		c.opHandler.Add(ops.Cs(lenReg), ops.Cs(regIndexStart), ops.Cs(regIndexStart)),
	)
	cmd += c.opHandler.ExecCond(
		fmt.Sprintf("score %s %s matches ..-1", ops.Cs(regIndexEnd), c.Namespace),
		true,
		c.opHandler.Add(ops.Cs(lenReg), ops.Cs(regIndexEnd), ops.Cs(regIndexEnd)),
	)

	// Move data to scoreboards
	cmd += c.opHandler.MoveScore(ops.Cs(lenReg), ops.Cs(lenReg))
	cmd += c.opHandler.MoveScore(ops.Cs(regIndexStart), ops.Cs(regIndexStart))
	cmd += c.opHandler.MoveScore(ops.Cs(regIndexEnd), ops.Cs(regIndexEnd))

	// Check if the start index is greater than the end index
	cmd += c.opHandler.ExecCond(
		fmt.Sprintf("score %s %s > %s %s", ops.Cs(regIndexStart), c.Namespace, ops.Cs(regIndexEnd), c.Namespace),
		true,
		c.opHandler.Exception("Start index greater than end index"),
	)
	switch expr.TargetExpr.ReturnType() {
	case expressions.StringType:
		cmd += c.opHandler.ExecCond(
			fmt.Sprintf("score %s %s >= %s %s", ops.Cs(regIndexStart), c.Namespace, ops.Cs(lenReg), c.Namespace),
			true,
			c.opHandler.Exception("Start slice index out of bounds"),
		)
		cmd += c.opHandler.ExecCond(
			fmt.Sprintf("score %s %s > %s %s", ops.Cs(regIndexEnd), c.Namespace, ops.Cs(lenReg), c.Namespace),
			true,
			c.opHandler.Exception("End slice index out of bounds"),
		)
	case expressions.ListIntType:
		fallthrough
	case expressions.ListDoubleType:
		fallthrough
	case expressions.ListStringType:
		cmd += c.opHandler.ExecCond(
			fmt.Sprintf("score %s %s >= %s %s", ops.Cs(regIndexStart), c.Namespace, ops.Cs(lenReg), c.Namespace),
			true,
			c.opHandler.Exception("Index out of bounds"),
		)
	}

	switch expr.TargetExpr.ReturnType() {
	case expressions.StringType:
		// String slice operation
		cmd += c.opHandler.SliceString(ops.Cs(ops.RX), ops.Cs(regIndexStart), ops.Cs(regIndexEnd), ops.Cs(ops.RX))
	case expressions.ListIntType:
		fallthrough
	case expressions.ListDoubleType:
		fallthrough
	case expressions.ListStringType:
		// List index operation
		if expr.EndIndex != nil {
			c.error(expr.SourceLocation, "List slicing is not supported")
			return ""
		}
		cmd += c.opHandler.GetListIndex(ops.Cs(ops.RX), ops.Cs(regIndexStart), ops.Cs(ops.RX))
	}
	return cmd
}

func (c *Compiler) VisitList(expr expressions.ListExpr) interface{} {
	cmd := ""
	cmd += "### List operation ###\n"
	regList := ops.Cs(c.newRegister(ops.RX))
	cmd += c.opHandler.MakeList(regList)
	for _, elem := range expr.Elements {
		cmd += elem.Accept(c).(string)
		cmd += c.opHandler.AppendList(regList, ops.Cs(ops.RX))
	}
	cmd += c.opHandler.Move(regList, ops.Cs(ops.RX))
	cmd += "### List operation end ###\n"
	return cmd
}
