package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	"path"
)

func (c *Compiler) VisitBinary(b *ast.BinaryExpr) any {
	regRa := c.makeReg(RA)
	regRb := c.makeReg(RB)

	cmd := c.n()

	cmd.Extend(ast.AcceptExpr[interfaces.IRCode](b.Left, c))
	cmd.CopyVar(RX, regRa)
	cmd.Extend(ast.AcceptExpr[interfaces.IRCode](b.Right, c))
	cmd.CopyVar(RX, regRb)

	switch b.Operator.Type {
	case tokens.EqualEqual, tokens.BangEqual, tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual:
		switch b.GetResolvedType() {
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
				c.error(b, "Invalid operator for type 'str'")
			}
		}
		return cmd
	default:
		// nothing
	}
	switch b.GetResolvedType() {
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
			c.error(b, fmt.Sprintf("Invalid operator for type 'int': %s", b.Operator.Lexeme))
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
			c.error(b, "Invalid operator for type 'double'")
		}
	case types.StringType:
		if b.Operator.Type == tokens.Plus {
			cmd.StringConcat(regRa, regRb, RX)
		} else {
			c.error(b, "Invalid operator for type 'str'")
		}
	}
	return cmd
}

func (c *Compiler) VisitGrouping(g *ast.GroupingExpr) any {
	return g.Expression.Accept(c)
}

func (c *Compiler) VisitLiteral(l *ast.LiteralExpr) any {
	return c.n().SetVar(RX, l.Value)
}

func (c *Compiler) VisitUnary(u *ast.UnaryExpr) any {
	cmd := c.n()
	switch u.ReturnType() {
	case types.IntType:
		switch u.Operator.Type {
		case tokens.Minus:
			zero := ast.NewLiteralExpr(nbt.NewInt(0), types.IntType, u.GetSourceLocation())
			tmp := &ast.BinaryExpr{
				Left: zero,
				Operator: tokens.Token{
					Type: tokens.Minus,
				},
				Right: u.Expression,
			}
			cmd.Extend(ast.AcceptExpr[interfaces.IRCode](tmp, c))
		case tokens.Bang:
			cmd.Extend(ast.AcceptExpr[interfaces.IRCode](u.Expression, c))
		default:
			c.error(u, "Invalid operator for type 'int'")
		}
	default:
		c.error(u, "Invalid type in unary expression")
	}
	return cmd
}

func (c *Compiler) VisitVariable(v *ast.VariableExpr) any {
	cmd := c.n()
	return cmd.CopyVar(v.Name.Lexeme, RX)
}

func (c *Compiler) VisitDotAccess(v *ast.DotAccessExpr) any {
	cmd := c.n()
	cmd.Extend(ast.AcceptExpr[interfaces.IRCode](v.Source, c))
	cmd.CopyVar(RX, RA)
	cmd.StructGet(RA, v.Name.Lexeme, RX)
	return cmd
}

func (c *Compiler) VisitCall(f *ast.CallExpr) any {
	cmd := c.n()
	ns, fn := utils.SplitFunctionName(f.GetResolvedName(), c.Namespace)
	sym, _ := c.currentScope.Lookup(f.GetResolvedName())
	switch sym.DeclarationNode().Type() {
	case ast.FunctionDeclarationStatement:
		funcStmt := sym.DeclarationNode().(ast.FunctionDeclarationStmt)
		for j, arg := range f.Arguments {
			cmd.Extend(ast.AcceptExpr[interfaces.IRCode](arg, c))
			argName := funcStmt.Parameters[j].Name.Lexeme
			cmd.CopyArg(RX, fn, argName)
		}
		if ns == c.Namespace {
			funcName := fmt.Sprintf("%s:%s", ns, path.Join(paths.FunctionBranches, fn))
			cmd.CallWithArgs(funcName, fmt.Sprintf("%s.%s", ArgPath, fn)) // Call wrapped function
		} else {
			cmd.Call(fmt.Sprintf("%s:%s", ns, fn))
		}
		if f.GetResolvedType() != types.VoidType {
			cmd.CopyVar(RET, RX)
		}
	case ast.StructDeclarationStatement:
		structStmt := sym.DeclarationNode().(ast.StructDeclarationStmt)
		structReg := c.makeReg(RX)
		cmd.SetVar(structReg, structStmt.StructType.ToNBT())
		for j, arg := range f.Arguments {
			cmd.Extend(ast.AcceptExpr[interfaces.IRCode](arg, c))
			fieldName := structStmt.StructType.GetFieldNames()[j]
			cmd.StructSet(RX, fieldName, structReg)
		}
		cmd.CopyVar(structReg, RX)
	}
	return cmd
}

func (c *Compiler) VisitLogical(l *ast.LogicalExpr) any {
	cmd := c.n()
	regRa := c.makeReg(RA)
	regRb := c.makeReg(RB)
	cmd.Extend(ast.AcceptExpr[interfaces.IRCode](l.Left, c))
	cmd.CopyVar(RX, regRa)
	cmd.Extend(ast.AcceptExpr[interfaces.IRCode](l.Right, c))
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
		c.error(l, fmt.Sprintf("Invalid operator: '%s'", l.Operator.Lexeme))
	}
	return cmd
}

func (c *Compiler) VisitSlice(s *ast.SliceExpr) any {
	cmd := c.n()
	regIndexStart := c.makeReg(RA)
	regIndexEnd := c.makeReg(RB)

	cmd.Extend(ast.AcceptExpr[interfaces.IRCode](s.StartIndex, c))
	cmd.CopyVar(RX, regIndexStart)
	if s.EndIndex == nil {
		cmd.CopyVar(RX, regIndexEnd)
	} else {
		cmd.Extend(ast.AcceptExpr[interfaces.IRCode](s.EndIndex, c))
		cmd.CopyVar(RX, regIndexEnd)
	}
	targetReg := c.makeReg(RX)
	cmd.Extend(ast.AcceptExpr[interfaces.IRCode](s.TargetExpr, c))
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
	if s.EndIndex != nil {
		cmd.IntCompare(regIndexStart, regIndexEnd, tokens.Greater, RX)
		cmd.Load(RX, RX)
		cmd.If(RX, c.n().ExceptionFormat(
			nbt.NewFormattedString("Exception at ", nbt.StringFormat{Color: nbt.Red, Italic: true}),
			nbt.NewFormattedString(s.EndIndex.GetSourceLocation().ToString(), nbt.StringFormat{Color: nbt.Gold}),
			nbt.NewFormattedString(": Invalid slice range. End index can't be smaller than start index", nbt.StringFormat{Color: nbt.Red, Italic: true})),
		)
		cmd.If(RX, c.n().Ret())
	}

	switch s.TargetExpr.GetResolvedType().(type) {
	case types.PrimitiveTypeStruct:
		switch s.TargetExpr.GetResolvedType() {
		case types.StringType:
			cmd.IntCompare(regIndexStart, lenReg, tokens.GreaterEqual, RX)
			cmd.Load(RX, RX)
			cmd.If(RX, c.n().ExceptionFormat(
				nbt.NewErrorString("OutOfBoundsException at "),
				nbt.NewFormattedString(s.EndIndex.GetSourceLocation().ToString(), nbt.StringFormat{Color: nbt.Gold}),
				nbt.NewErrorString(": Invalid slice range. Start index out of bounds"),
			))
			cmd.If(RX, c.n().Ret())
			if s.EndIndex != nil {
				cmd.IntCompare(regIndexEnd, lenReg, tokens.GreaterEqual, RX)
				cmd.Load(RX, RX)
				cmd.If(RX, c.n().ExceptionFormat(
					nbt.NewErrorString("OutOfBoundsException at ", nbt.StringFormat{Color: nbt.Red, Italic: true}),
					nbt.NewFormattedString(s.EndIndex.GetSourceLocation().ToString(), nbt.StringFormat{Color: nbt.Gold}),
					nbt.NewErrorString(": Invalid slice range. End index out of bounds"),
				))
				cmd.If(RX, c.n().Ret())
			}
			cmd.StringSlice(targetReg, regIndexStart, regIndexEnd, RX)
		}
	case types.ListTypeStruct:
		cmd.IntCompare(regIndexStart, lenReg, tokens.GreaterEqual, RX)
		cmd.Load(RX, RX)
		cmd.If(RX, c.n().ExceptionFormat(
			nbt.NewErrorString("OutOfBoundsException at "),
			nbt.NewFormattedString(s.StartIndex.GetSourceLocation().ToString(), nbt.StringFormat{Color: nbt.Gold}),
			nbt.NewErrorString(": Index out of bounds"),
		))
		cmd.If(RX, c.n().Ret())
		if s.EndIndex != nil {
			c.error(s, "List slices are not supported.")
		}
		cmd.MakeIndex(regIndexStart, regIndexStart)
		cmd.PathGet(targetReg, regIndexStart, RX)
	}
	return cmd
}

func (c *Compiler) VisitList(s *ast.ListExpr) any {
	cmd := c.n()
	regList := c.makeReg(RX)
	cmd.SetVar(regList, nbt.NewList())
	for _, elem := range s.Elements {
		cmd.Extend(ast.AcceptExpr[interfaces.IRCode](elem, c))
		cmd.AppendCopy(regList, RX)
	}
	cmd.CopyVar(regList, RX)
	return cmd
}
