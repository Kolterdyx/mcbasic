package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/compiler/mapping"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"strconv"
)

func (c *Compiler) VisitLiteral(expr expressions.LiteralExpr) string {
	switch expr.ReturnType() {
	case types.IntType:
		i, err := strconv.ParseInt(expr.Value, 10, 64)
		if err != nil {
			c.error(expr.SourceLocation, "Invalid integer literal")
		}
		return c.commandMapper.MakeConst(nbt.NewInt(i), c.commandMapper.Cs(mapping.RX))
	case types.StringType:
		return c.commandMapper.MakeConst(nbt.NewString(expr.Value), c.commandMapper.Cs(mapping.RX))
	case types.DoubleType:
		d, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			c.error(expr.SourceLocation, "Invalid double literal")
		}
		return c.commandMapper.MakeConst(nbt.NewDouble(d), c.commandMapper.Cs(mapping.RX))
	default:
		c.error(expr.SourceLocation, "Invalid type in literal expression")
	}
	return ""
}

func (c *Compiler) VisitBinary(expr expressions.BinaryExpr) string {
	cmd := ""

	cmd, regRa := c.commandMapper.MakeRegister(mapping.RA)
	cmd, regRb := c.commandMapper.MakeRegister(mapping.RB)

	cmd += "### BEGIN Binary operation ###\n"
	cmd += "###       Left side ###\n"
	cmd += expr.Left.Accept(c)
	cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(regRa))
	cmd += "###       Right side ###\n"
	cmd += expr.Right.Accept(c)
	cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(regRb))

	switch expr.Operator.Type {
	case tokens.EqualEqual, tokens.BangEqual, tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual:
		cmd += c.Compare(expr, c.commandMapper.Cs(regRa), c.commandMapper.Cs(regRb), c.commandMapper.Cs(mapping.RX))
		return cmd
	default:
		// Do nothing
	}

	switch expr.ReturnType() {
	case types.IntType:
		switch expr.Operator.Type {
		case tokens.Plus:
			cmd += c.commandMapper.IntAdd(regRa, regRb, mapping.RX)
		case tokens.Minus:
			cmd += c.commandMapper.IntSub(regRa, regRb, mapping.RX)
		case tokens.Star:
			cmd += c.commandMapper.IntMul(regRa, regRb, mapping.RX)
		case tokens.Slash:
			cmd += c.commandMapper.IntDiv(regRa, regRb, mapping.RX)
		case tokens.Percent:
			cmd += c.commandMapper.IntMod(regRa, regRb, mapping.RX)
		default:
			c.error(expr.SourceLocation, "Invalid operator for integers")
		}
	case types.DoubleType:
		switch expr.Operator.Type {
		case tokens.Plus:
			cmd += c.commandMapper.DoubleAdd(regRa, regRb, mapping.RX)
		case tokens.Minus:
			cmd += c.commandMapper.DoubleSub(regRa, regRb, mapping.RX)
		case tokens.Star:
			cmd += c.commandMapper.DoubleMul(regRa, regRb, mapping.RX)
		case tokens.Slash:
			cmd += c.commandMapper.DoubleDiv(regRa, regRb, mapping.RX)
		case tokens.Percent:
			cmd += c.commandMapper.DoubleMod(regRa, regRb, mapping.RX)
		default:
			c.error(expr.SourceLocation, "Invalid operator for double numbers")
		}

	case types.StringType:
		if expr.Operator.Type == tokens.Plus {
			cmd += c.commandMapper.Concat(c.commandMapper.Cs(regRa), c.commandMapper.Cs(regRb), c.commandMapper.Cs(mapping.RX))
		} else {
			c.error(expr.SourceLocation, "Invalid operator for strings")
		}
	default:
		c.error(expr.SourceLocation, "Invalid type combination in binary operation")
	}
	cmd += "### END   Binary operation ###\n"
	return cmd
}

func (c *Compiler) VisitVariable(expr expressions.VariableExpr) string {
	return c.commandMapper.Move(c.commandMapper.Cs(expr.Name.Lexeme), c.commandMapper.Cs(mapping.RX))
}

func (c *Compiler) VisitFieldAccess(expr expressions.FieldAccessExpr) string {
	cmd := "### BEGIN Struct field access operation ###\n"
	cmd += expr.Source.Accept(c)
	cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(mapping.RA))
	cmd += c.commandMapper.StructGet(c.commandMapper.Cs(mapping.RA), expr.Field.Lexeme, c.commandMapper.Cs(mapping.RX))
	cmd += "### END   Struct field access operation ###\n"
	return cmd
}

func (c *Compiler) VisitFunctionCall(expr expressions.FunctionCallExpr) string {
	cmd := ""
	for i, arg := range expr.Arguments {
		cmd += arg.Accept(c)
		argName := c.functions[expr.Name.Lexeme].Args[i].Name
		cmd += c.commandMapper.LoadArg(expr.Name.Lexeme, argName, c.commandMapper.Cs(mapping.RX))
	}
	if expr.ReturnType() != types.VoidType {
		cmd += c.commandMapper.Call(expr.Name.Lexeme, c.commandMapper.Cs(mapping.RX))
	} else {
		cmd += c.commandMapper.Call(expr.Name.Lexeme, "")
	}
	return cmd
}

func (c *Compiler) VisitUnary(expr expressions.UnaryExpr) string {
	cmd := ""
	switch expr.ReturnType() {
	case types.IntType:
		switch expr.Operator.Type {
		case tokens.Minus:
			zero := expressions.LiteralExpr{Value: "0", SourceLocation: expr.SourceLocation, ValueType: types.IntType}
			tmp := expressions.BinaryExpr{
				Left: zero,
				Operator: tokens.Token{
					Type: tokens.Minus,
				},
				Right: expr.Expression,
			}
			cmd += tmp.Accept(c)
		case tokens.Bang:
			cmd += expr.Expression.Accept(c)
			cmd += c.commandMapper.NegateNumber(c.commandMapper.Cs(mapping.RX))
		default:
			c.error(expr.SourceLocation, "Invalid operator for integers")
		}
	default:
		c.error(expr.SourceLocation, "Invalid type in unary expression")
	}
	return cmd
}

func (c *Compiler) VisitGrouping(expr expressions.GroupingExpr) string {
	return expr.Expression.Accept(c)
}

func (c *Compiler) VisitLogical(expr expressions.LogicalExpr) string {
	leftSide := ""
	rightSide := ""

	cmd := ""
	regRa := c.newRegister(mapping.RA)
	cmd += c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(regRa))
	regRb := c.newRegister(mapping.RB)
	cmd += c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(regRb))

	cmd += "### BEGIN Logical operation ###\n"
	leftSide += "###       Logical operation left side ###\n"
	leftSide += expr.Left.Accept(c)
	leftSide += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), regRa)
	rightSide += "###       Logical operation right side ###\n"
	rightSide += expr.Right.Accept(c)
	rightSide += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), regRb)
	rightSide += c.commandMapper.MoveScore(regRb, regRb)

	cmd += leftSide
	switch expr.Operator.Type {
	case tokens.And:
		// If left side is false, return false
		evalRightSide := ""
		cmd += c.commandMapper.MoveScore(regRa, regRa)
		cmd += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), true, c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(mapping.RX)))
		evalRightSide += rightSide
		evalRightSide += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), true, c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(mapping.RX)))
		evalRightSide += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), false, c.commandMapper.Move(regRb, c.commandMapper.Cs(mapping.RX)))
		cmd += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), false, evalRightSide)
	case tokens.Or:
		// If left side is true, return true
		evalRightSide := ""
		cmd += c.commandMapper.MoveScore(regRa, regRa)
		cmd += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), false, c.commandMapper.Move(regRa, c.commandMapper.Cs(mapping.RX)))
		evalRightSide += rightSide
		evalRightSide += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), true, c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(mapping.RX)))
		evalRightSide += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), false, c.commandMapper.Move(regRb, c.commandMapper.Cs(mapping.RX)))
		cmd += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), true, evalRightSide)
	default:
		c.error(expr.SourceLocation, "Invalid operator for logical expressions")
	}

	cmd += "### END   Logical operation ###\n"
	return cmd
}

func (c *Compiler) VisitSlice(expr expressions.SliceExpr) string {
	regIndexStart := c.newRegister(mapping.RA)
	regIndexEnd := c.newRegister(mapping.RB)

	cmd := ""
	cmd += c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(regIndexStart))
	cmd += c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(regIndexEnd))
	cmd += "### BEGIN String slice operation ###\n"
	cmd += "###       Accept start index ###\n"
	cmd += expr.StartIndex.Accept(c)
	cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(regIndexStart))

	if expr.EndIndex == nil {
		cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(regIndexEnd))
		cmd += c.commandMapper.IncScore(c.commandMapper.Cs(regIndexEnd))
	} else {
		cmd += "###       Accept end index ###\n"
		cmd += expr.EndIndex.Accept(c)
		cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(regIndexEnd))
	}
	cmd += expr.TargetExpr.Accept(c)

	// Check index bounds

	cmd += "###       Cheking index bounds ###\n"
	lenReg := c.newRegister(mapping.RX)
	cmd += c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(lenReg))

	cmd += c.commandMapper.Size(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(lenReg))

	// If any of the indexes are negative, add the length of the string to them
	cmd += c.commandMapper.ExecCond(
		fmt.Sprintf("score %s %s matches ..-1", c.commandMapper.Cs(regIndexStart), c.Namespace),
		true,
		c.commandMapper.IntAdd(c.commandMapper.Cs(lenReg), c.commandMapper.Cs(regIndexStart), c.commandMapper.Cs(regIndexStart)),
	)
	cmd += c.commandMapper.ExecCond(
		fmt.Sprintf("score %s %s matches ..-1", c.commandMapper.Cs(regIndexEnd), c.Namespace),
		true,
		c.commandMapper.IntAdd(c.commandMapper.Cs(lenReg), c.commandMapper.Cs(regIndexEnd), c.commandMapper.Cs(regIndexEnd)),
	)

	// Move data to scoreboards
	cmd += c.commandMapper.MoveScore(c.commandMapper.Cs(lenReg), c.commandMapper.Cs(lenReg))
	cmd += c.commandMapper.MoveScore(c.commandMapper.Cs(regIndexStart), c.commandMapper.Cs(regIndexStart))
	cmd += c.commandMapper.MoveScore(c.commandMapper.Cs(regIndexEnd), c.commandMapper.Cs(regIndexEnd))

	// Check if the start index is greater than the end index
	cmd += c.commandMapper.ExecCond(
		fmt.Sprintf("score %s %s > %s %s", c.commandMapper.Cs(regIndexStart), c.Namespace, c.commandMapper.Cs(regIndexEnd), c.Namespace),
		true,
		c.commandMapper.Exception("Start index greater than end index"),
	)
	switch expr.TargetExpr.ReturnType().(type) {
	case types.PrimitiveTypeStruct:
		switch expr.TargetExpr.ReturnType() {
		case types.StringType:
			cmd += c.commandMapper.ExecCond(
				fmt.Sprintf("score %s %s >= %s %s", c.commandMapper.Cs(regIndexStart), c.Namespace, c.commandMapper.Cs(lenReg), c.Namespace),
				true,
				c.commandMapper.Exception("Start slice index out of bounds"),
			)
			cmd += c.commandMapper.ExecCond(
				fmt.Sprintf("score %s %s > %s %s", c.commandMapper.Cs(regIndexEnd), c.Namespace, c.commandMapper.Cs(lenReg), c.Namespace),
				true,
				c.commandMapper.Exception("End slice index out of bounds"),
			)
			cmd += "###       Slice string ###\n"
			cmd += c.commandMapper.SliceString(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(regIndexStart), c.commandMapper.Cs(regIndexEnd), c.commandMapper.Cs(mapping.RX))
		}
	case types.ListTypeStruct:
		cmd += c.commandMapper.ExecCond(
			fmt.Sprintf("score %s %s >= %s %s", c.commandMapper.Cs(regIndexStart), c.Namespace, c.commandMapper.Cs(lenReg), c.Namespace),
			true,
			c.commandMapper.Exception("Index out of bounds"),
		)
		cmd += "###       Index list ###\n"
		if expr.EndIndex != nil {
			c.error(expr.SourceLocation, "List slicing is not supported")
			return ""
		}
		cmd += c.commandMapper.MakeIndex(c.commandMapper.Cs(regIndexStart), c.commandMapper.Cs(regIndexStart))
		cmd += c.commandMapper.PathGet(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(regIndexStart), c.commandMapper.Cs(mapping.RX))
	}
	cmd += "### END   String slice operation ###\n"
	return cmd
}

func (c *Compiler) VisitList(expr expressions.ListExpr) string {
	cmd := "### BEGIN List init operation ###\n"
	regList := c.commandMapper.Cs(c.newRegister(mapping.RX))
	cmd += c.commandMapper.MakeConst(nbt.NewInt(0), regList)
	cmd += c.commandMapper.MakeConst(nbt.NewList(), regList)
	for _, elem := range expr.Elements {
		cmd += elem.Accept(c)
		cmd += c.commandMapper.AppendList(regList, c.commandMapper.Cs(mapping.RX))
	}
	cmd += c.commandMapper.Move(regList, c.commandMapper.Cs(mapping.RX))
	cmd += "### END   List operation ###\n"
	return cmd
}

func (c *Compiler) VisitStruct(expr expressions.StructExpr) string {
	cmd := "### BEGIN Struct init operation ###\n"
	regStruct := c.commandMapper.Cs(c.newRegister(mapping.RX))
	cmd += c.commandMapper.MakeConst(expr.StructType.ToNBT(), regStruct)
	structFields := expr.StructType.GetFieldNames()
	for i, arg := range expr.Args {
		cmd += arg.Accept(c)
		cmd += c.commandMapper.StructSet(c.commandMapper.Cs(mapping.RX), structFields[i], regStruct)
	}
	cmd += c.commandMapper.Move(regStruct, c.commandMapper.Cs(mapping.RX))
	cmd += "### END   Struct operation ###\n"
	return cmd
}
