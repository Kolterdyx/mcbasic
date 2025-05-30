package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	"path"
)

func (c *Compiler) VisitBinary(b expressions.BinaryExpr) interfaces.IRCode {
	regRa := c.makeReg(RA)
	regRb := c.makeReg(RB)

	cmd := c.n()

	cmd.Extend(b.Left.Accept(c))
	cmd.CopyVar(RX, regRa)
	cmd.Extend(b.Right.Accept(c))
	cmd.CopyVar(RX, regRb)

	switch b.Operator.Type {
	case tokens.EqualEqual, tokens.BangEqual, tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual:
		switch b.ReturnType() {
		case types.IntType:
			cmd.Load(regRa, regRa)
			cmd.Load(regRb, regRb)
			cmd.IntCompare(regRa, regRb, b.Operator.Type, RX)
		case types.DoubleType:
			cmd.DoubleCompare(regRa, regRb, b.Operator.Type, RX)
		case types.StringType:
			if b.Operator.Type == tokens.EqualEqual {
				cmd.StringCompare(regRa, regRb, RX)
			} else {
				c.error(b.SourceLocation, "Invalid operator for type 'str'")
			}
		}
		return cmd
	default:
		// nothing
	}
	switch b.ReturnType() {
	case types.IntType:
		switch b.Operator.Type {
		case tokens.Plus:
			cmd.IntAdd(regRa, regRb, RX)
		case tokens.Minus:
			cmd.IntSub(regRa, regRb, RX)
		case tokens.Star:
			cmd.IntMul(regRa, regRb, RX)
		case tokens.Slash:
			cmd.IntDiv(regRa, regRb, RX)
		case tokens.Percent:
			cmd.IntMod(regRa, regRb, RX)
		default:
			c.error(b.SourceLocation, fmt.Sprintf("Invalid operator for type 'int': %s", b.Operator.Lexeme))
		}
	case types.DoubleType:
		switch b.Operator.Type {
		case tokens.Plus:
			cmd.DoubleAdd(regRa, regRb, RX)
		case tokens.Minus:
			cmd.DoubleSub(regRa, regRb, RX)
		case tokens.Star:
			cmd.DoubleMul(regRa, regRb, RX)
		case tokens.Slash:
			cmd.DoubleDiv(regRa, regRb, RX)
		case tokens.Percent:
			cmd.DoubleMod(regRa, regRb, RX)
		default:
			c.error(b.SourceLocation, "Invalid operator for type 'double'")
		}
	case types.StringType:
		if b.Operator.Type == tokens.Plus {
			cmd.StringConcat(regRa, regRb, RX)
		} else {
			c.error(b.SourceLocation, "Invalid operator for type 'str'")
		}
	}
	return cmd
}

func (c *Compiler) VisitGrouping(g expressions.GroupingExpr) interfaces.IRCode {
	return g.Expression.Accept(c)
}

func (c *Compiler) VisitLiteral(l expressions.LiteralExpr) interfaces.IRCode {
	return c.n().SetVar(RX, l.Value)
}

func (c *Compiler) VisitUnary(u expressions.UnaryExpr) interfaces.IRCode {
	cmd := c.n()
	switch u.ReturnType() {
	case types.IntType:
		switch u.Operator.Type {
		case tokens.Minus:
			zero := expressions.LiteralExpr{Value: nbt.NewInt(0), SourceLocation: u.SourceLocation, ValueType: types.IntType}
			tmp := expressions.BinaryExpr{
				Left: zero,
				Operator: tokens.Token{
					Type: tokens.Minus,
				},
				Right: u.Expression,
			}
			cmd.Extend(tmp.Accept(c))
		case tokens.Bang:
			cmd.Extend(u.Expression.Accept(c))
		default:
			c.error(u.SourceLocation, "Invalid operator for type 'int'")
		}
	default:
		c.error(u.SourceLocation, "Invalid type in unary expression")
	}
	return cmd
}

func (c *Compiler) VisitVariable(v expressions.VariableExpr) interfaces.IRCode {
	cmd := c.n()
	return cmd.CopyVar(v.Name.Lexeme, RX)
}

func (c *Compiler) VisitFieldAccess(v expressions.FieldAccessExpr) interfaces.IRCode {
	cmd := c.n()
	cmd.Extend(v.Source.Accept(c))
	cmd.CopyVar(RX, RA)
	cmd.StructGet(RA, v.Field.Lexeme, RX)
	return cmd
}

func (c *Compiler) VisitFunctionCall(f expressions.FunctionCallExpr) interfaces.IRCode {
	cmd := c.n()
	ns, fn := utils.SplitFunctionName(f.Name.Lexeme, c.Namespace)
	for j, arg := range f.Arguments {
		cmd.Extend(arg.Accept(c))
		argName := c.functionDefinitions[f.Name.Lexeme].Args[j].Name
		cmd.CopyArg(RX, fn, argName)
	}
	if ns == c.Namespace {
		funcName := fmt.Sprintf("%s:%s", ns, path.Join(paths.FunctionBranches, fn))
		cmd.CallWithArgs(funcName, fmt.Sprintf("%s.%s", ArgPath, fn)) // Call wrapped function
	} else {
		cmd.Call(fmt.Sprintf("%s:%s", ns, fn))
	}
	if f.ReturnType() != types.VoidType {
		cmd.CopyVar(RET, RX)
	}
	return cmd
}

func (c *Compiler) VisitLogical(l expressions.LogicalExpr) interfaces.IRCode {
	cmd := c.n()
	regRa := c.makeReg(RA)
	regRb := c.makeReg(RB)
	cmd.Extend(l.Left.Accept(c))
	cmd.CopyVar(RX, regRa)
	cmd.Extend(l.Right.Accept(c))
	cmd.CopyVar(RX, regRb)

	switch l.Operator.Type {
	case tokens.And:
		cmd.Load(regRa, regRa)
		cmd.If(regRa, c.n().SetVar(regRa, nbt.NewInt(1)))
		cmd.If(regRa, c.n().If(regRb, c.n().SetVar(RX, nbt.NewInt(1))))
	case tokens.Or:
		cmd.Load(regRa, regRa)
		cmd.If(regRa, c.n().SetVar(RX, nbt.NewInt(1)))
		cmd.Load(regRb, regRb)
		cmd.If(regRb, c.n().SetVar(RX, nbt.NewInt(1)))
	default:
		c.error(l.SourceLocation, fmt.Sprintf("Invalid operator: '%s'", l.Operator.Lexeme))
	}
	return cmd
}

func (c *Compiler) VisitSlice(s expressions.SliceExpr) interfaces.IRCode {
	cmd := c.n()
	regIndexStart := c.makeReg(RA)
	regIndexEnd := c.makeReg(RB)

	cmd.Extend(s.StartIndex.Accept(c))
	cmd.CopyVar(RX, regIndexStart)
	if s.EndIndex == nil {
		cmd.CopyVar(RX, regIndexEnd)
	} else {
		cmd.Extend(s.EndIndex.Accept(c))
		cmd.CopyVar(RX, regIndexEnd)
	}
	targetReg := c.makeReg(RX)
	cmd.Extend(s.TargetExpr.Accept(c))
	cmd.CopyVar(RX, targetReg)
	lenReg := c.makeReg(RX)
	cmd.Size(targetReg, lenReg)
	cmd.Load(lenReg, lenReg)
	cmd.Score(RX, nbt.NewInt(-1))
	cmd.Load(regIndexStart, regIndexStart)
	cmd.IntCompare(regIndexStart, RX, tokens.LessEqual, RX)
	cmd.Load(RX, RX)
	cmd.If(RX, c.n().IntAdd(regIndexStart, lenReg, regIndexStart))
	cmd.If(RX, c.n().Load(regIndexStart, regIndexStart))
	cmd.Load(regIndexEnd, regIndexEnd)
	cmd.IntCompare(regIndexEnd, RX, tokens.LessEqual, RX)
	cmd.Load(RX, RX)
	cmd.If(RX, c.n().IntAdd(regIndexEnd, lenReg, regIndexEnd))
	cmd.If(RX, c.n().Load(regIndexEnd, regIndexEnd))
	if s.EndIndex == nil {
		cmd.IntCompare(regIndexStart, regIndexEnd, tokens.Greater, RX)
		cmd.Load(RX, RX)
		cmd.If(RX, c.n().Exception(fmt.Sprintf("Exception at %s: Invalid slice range. End index can't be smaller than start index", s.SourceLocation.ToString())))
		cmd.If(RX, c.n().Ret())
	}

	switch s.TargetExpr.ReturnType().(type) {
	case types.PrimitiveTypeStruct:
		switch s.TargetExpr.ReturnType() {
		case types.StringType:
			cmd.IntCompare(regIndexStart, lenReg, tokens.GreaterEqual, RX)
			cmd.Load(RX, RX)
			cmd.If(RX, c.n().Exception(fmt.Sprintf("Exception at %s: Invalid slice range. Start index out of bounds", s.SourceLocation.ToString())))
			cmd.If(RX, c.n().Ret())
			if s.EndIndex != nil {
				cmd.IntCompare(regIndexEnd, lenReg, tokens.GreaterEqual, RX)
				cmd.Load(RX, RX)
				cmd.If(RX, c.n().Exception(fmt.Sprintf("Exception at %s: Invalid slice range. Start index out of bounds", s.SourceLocation.ToString())))
				cmd.If(RX, c.n().Ret())
			}
			cmd.StringSlice(targetReg, regIndexStart, regIndexEnd, RX)
		}
	case types.ListTypeStruct:
		cmd.IntCompare(regIndexStart, lenReg, tokens.GreaterEqual, RX)
		cmd.Load(RX, RX)
		cmd.If(RX, c.n().Exception(fmt.Sprintf("Exception at %s: Index out of bounds", s.SourceLocation.ToString())))
		cmd.If(RX, c.n().Ret())
		if s.EndIndex != nil {
			c.error(s.SourceLocation, "List slices are not supported.")
		}
		cmd.MakeIndex(regIndexStart, regIndexStart)
		cmd.PathGet(targetReg, regIndexStart, RX)
	}
	return cmd
}

func (c *Compiler) VisitList(s expressions.ListExpr) interfaces.IRCode {
	cmd := c.n()
	regList := c.makeReg(RX)
	cmd.SetVar(regList, nbt.NewList())
	for _, elem := range s.Elements {
		cmd.Extend(elem.Accept(c))
		cmd.AppendCopy(regList, RX)
	}
	cmd.CopyVar(regList, RX)
	return cmd
}

func (c *Compiler) VisitStruct(s expressions.StructExpr) interfaces.IRCode {
	cmd := c.n()
	regStruct := c.makeReg(RX)
	cmd.SetVar(regStruct, s.StructType.ToNBT())
	fieldNames := s.StructType.GetFieldNames()
	for j, args := range s.Args {
		cmd.Extend(args.Accept(c))
		cmd.StructSet(c.varPath(RX), fieldNames[j], regStruct)
	}
	cmd.CopyVar(regStruct, RX)
	return cmd
}
