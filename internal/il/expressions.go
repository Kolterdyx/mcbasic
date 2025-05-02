package il

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"strconv"
)

func (c *Compiler) VisitBinary(b expressions.BinaryExpr) (cmd string) {
	regRa := c.makeReg(RA)
	regRb := c.makeReg(RB)
	cmd += b.Left.Accept(c)
	cmd += c.CopyVar(RX, regRa)
	cmd += b.Right.Accept(c)
	cmd += c.CopyVar(RX, regRb)

	switch b.Operator.Type {
	case tokens.EqualEqual, tokens.BangEqual, tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual:
		switch b.ReturnType() {
		case types.IntType:
			cmd += c.IntCompare(regRa, regRb, b.Operator.Type, RX)
		case types.DoubleType:
			cmd += c.DoubleCompare(regRa, regRb, b.Operator.Type, RX)
		case types.StringType:
			if b.Operator.Type == tokens.EqualEqual {
				cmd += c.StringCompare(regRa, regRb, RX)
			} else {
				c.error(b.SourceLocation, "Invalid operator for type 'str'")
			}
		}
		return
	default:
		// nothing
	}
	switch b.ReturnType() {
	case types.IntType:
		switch b.Operator.Type {
		case tokens.Plus:
			cmd += c.IntAdd(regRa, regRb, RX)
		case tokens.Minus:
			cmd += c.IntSub(regRa, regRb, RX)
		case tokens.Star:
			cmd += c.IntMul(regRa, regRb, RX)
		case tokens.Slash:
			cmd += c.IntDiv(regRa, regRb, RX)
		case tokens.Percent:
			cmd += c.IntMod(regRa, regRb, RX)
		default:
			c.error(b.SourceLocation, fmt.Sprintf("Invalid operator for type 'int': %s", b.Operator.Lexeme))
		}
	case types.DoubleType:
		switch b.Operator.Type {
		case tokens.Plus:
			cmd += c.DoubleAdd(regRa, regRb, RX)
		case tokens.Minus:
			cmd += c.DoubleSub(regRa, regRb, RX)
		case tokens.Star:
			cmd += c.DoubleMul(regRa, regRb, RX)
		case tokens.Slash:
			cmd += c.DoubleDiv(regRa, regRb, RX)
		case tokens.Percent:
			cmd += c.DoubleMod(regRa, regRb, RX)
		default:
			c.error(b.SourceLocation, "Invalid operator for type 'double'")
		}
	case types.StringType:
		if b.Operator.Type == tokens.Plus {
			cmd += c.StringConcat(regRa, regRb, RX)
		} else {
			c.error(b.SourceLocation, "Invalid operator for type 'str'")
		}
	}
	return
}

func (c *Compiler) VisitGrouping(g expressions.GroupingExpr) (cmd string) {
	return g.Expression.Accept(c)
}

func (c *Compiler) VisitLiteral(l expressions.LiteralExpr) (cmd string) {
	switch l.ReturnType() {
	case types.StringType:
		return c.SetVar(RX, nbt.NewString(l.Value))
	case types.IntType:
		integer, err := strconv.ParseInt(l.Value, 10, 64)
		if err != nil {
			c.error(l.SourceLocation, "Invalid integer literal")
		}
		return c.SetVar(RX, nbt.NewInt(integer))
	case types.DoubleType:
		double, err := strconv.ParseFloat(l.Value, 64)
		if err != nil {
			c.error(l.SourceLocation, "Invalid double literal")
		}
		return c.SetVar(RX, nbt.NewDouble(double))
	default:
		c.error(l.SourceLocation, "Invalid type in literal expression")
	}
	return
}

func (c *Compiler) VisitUnary(u expressions.UnaryExpr) (cmd string) {
	switch u.ReturnType() {
	case types.IntType:
		switch u.Operator.Type {
		case tokens.Minus:
			zero := expressions.LiteralExpr{Value: "0", SourceLocation: u.SourceLocation, ValueType: types.IntType}
			tmp := expressions.BinaryExpr{
				Left: zero,
				Operator: tokens.Token{
					Type: tokens.Minus,
				},
				Right: u.Expression,
			}
			cmd += tmp.Accept(c)
		case tokens.Bang:
			cmd += u.Expression.Accept(c)
		default:
			c.error(u.SourceLocation, "Invalid operator for type 'int'")
		}
	default:
		c.error(u.SourceLocation, "Invalid type in unary expression")
	}
	return
}

func (c *Compiler) VisitVariable(v expressions.VariableExpr) (cmd string) {
	return c.CopyVar(v.Name.Lexeme, RX)
}

func (c *Compiler) VisitFieldAccess(v expressions.FieldAccessExpr) (cmd string) {
	cmd += v.Source.Accept(c)
	cmd += c.CopyVar(RX, RA)
	cmd += c.StructGet(RA, v.Field.Lexeme, RX)
	return
}

func (c *Compiler) VisitFunctionCall(f expressions.FunctionCallExpr) (cmd string) {
	for j, arg := range f.Arguments {
		cmd += arg.Accept(c)
		argName := c.functions[f.Name.Lexeme].Args[j].Name
		cmd += c.Copy(c.varPath(RX), c.argPath(f.Name.Lexeme, argName))
	}
	cmd += c.Call(f.Name.Lexeme)
	if f.ReturnType() != types.VoidType {
		cmd += c.CopyVar(RET, RX)
	}
	return
}

func (c *Compiler) VisitLogical(l expressions.LogicalExpr) (cmd string) {
	regRa := c.makeReg(RA)
	regRb := c.makeReg(RB)
	cmd += l.Left.Accept(c)
	cmd += c.CopyVar(RX, regRa)
	cmd += l.Right.Accept(c)
	cmd += c.CopyVar(RX, regRb)

	switch l.Operator.Type {
	case tokens.And:
		cmd += c.If(regRa, c.SetVar(regRa, nbt.NewInt(1)))
		cmd += c.If(regRa, c.If(regRb, c.SetVar(RX, nbt.NewInt(1))))
	case tokens.Or:
		cmd += c.If(regRa, c.SetVar(RX, nbt.NewInt(1)))
		cmd += c.If(regRb, c.SetVar(RX, nbt.NewInt(1)))
	default:
		c.error(l.SourceLocation, fmt.Sprintf("Invalid operator: '%s'", l.Operator.Lexeme))
	}
	return
}

func (c *Compiler) VisitSlice(s expressions.SliceExpr) (cmd string) {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitList(s expressions.ListExpr) (cmd string) {
	regList := c.makeReg(RX)
	cmd += c.SetVar(regList, nbt.NewList())
	for _, elem := range s.Elements {
		cmd += elem.Accept(c)
		cmd += c.Append(regList, RX)
	}
	cmd += c.CopyVar(RX, regList)
	return
}

func (c *Compiler) VisitStruct(s expressions.StructExpr) (cmd string) {
	regStruct := c.makeReg(RX)
	cmd += c.SetVar(regStruct, s.StructType.ToNBT())
	fieldNames := s.StructType.GetFieldNames()
	for j, args := range s.Args {
		cmd += args.Accept(c)
		cmd += c.StructSet(c.varPath(RX), fieldNames[j], regStruct)
	}
	cmd += c.CopyVar(regStruct, RX)
	return
}
