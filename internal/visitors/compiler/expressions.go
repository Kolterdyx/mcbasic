package compiler

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func (c *Compiler) VisitLiteral(expr expressions.LiteralExpr) interface{} {
	if expr.ReturnType() == expressions.NumberType {
		return c.opHandler.MoveConst(expr.Value.(string), ops.Cs(ops.RX))
	} else if expr.ReturnType() == expressions.StringType {
		return c.opHandler.MoveConst(strconv.Quote(expr.Value.(string)), ops.Cs(ops.RX))
	} else {
		log.Fatalln("Invalid type in literal expression")
		return ""
	}
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

	if expr.Left.ReturnType() != expr.Right.ReturnType() {
		log.Fatalln("Different types in binary operation")
	}
	cmd += "### Binary operation ###\n"

	switch expr.Operator.Type {
	case tokens.EqualEqual:
		fallthrough
	case tokens.BangEqual:
		fallthrough
	case tokens.Greater:
		fallthrough
	case tokens.GreaterEqual:
		fallthrough
	case tokens.Less:
		fallthrough
	case tokens.LessEqual:
		cmd += c.Compare(expr, ops.Cs(regRa), ops.Cs(regRb), ops.Cs(ops.RX))
		return cmd
	}

	if expr.ReturnType() == expressions.NumberType {
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
		}
	} else if expr.ReturnType() == expressions.StringType {
		log.Fatalln("String operations are not supported yet")
		//if expr.Operator.Type == tokens.Plus {
		//	return c.opHandler.Concat(ops.RA, ops.RB, ops.RX)
		//} else {
		//	panic("Unknown operator")
		//}
	} else {
		log.Fatalln("Invalid type in binary operation")
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
