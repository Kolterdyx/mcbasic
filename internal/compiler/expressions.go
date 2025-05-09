package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/compiler/ops"
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
		return c.opHandler.MakeConst(nbt.NewInt(i), ops.Cs(ops.RX))
	case types.StringType:
		return c.opHandler.MakeConst(nbt.NewString(expr.Value), ops.Cs(ops.RX))
	case types.DoubleType:
		d, err := strconv.ParseFloat(expr.Value, 64)
		if err != nil {
			c.error(expr.SourceLocation, "Invalid double literal")
		}
		return c.opHandler.MakeConst(nbt.NewDouble(d), ops.Cs(ops.RX))
	default:
		c.error(expr.SourceLocation, "Invalid type in literal expression")
	}
	return ""
}

func (c *Compiler) VisitBinary(expr expressions.BinaryExpr) string {
	cmd := ""

	regRa := c.newRegister(ops.RA)
	cmd += c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(regRa))
	regRb := c.newRegister(ops.RB)
	cmd += c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(regRb))

	cmd += "### BEGIN Binary operation ###\n"
	cmd += "###       Left side ###\n"
	cmd += expr.Left.Accept(c)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regRa))
	cmd += "###       Right side ###\n"
	cmd += expr.Right.Accept(c)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regRb))

	switch expr.Operator.Type {
	case tokens.EqualEqual, tokens.BangEqual, tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual:
		cmd += c.Compare(expr, ops.Cs(regRa), ops.Cs(regRb), ops.Cs(ops.RX))
		return cmd
	default:
		// Do nothing
	}

	switch expr.ReturnType() {
	case types.IntType:
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
	case types.DoubleType:
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

	case types.StringType:
		if expr.Operator.Type == tokens.Plus {
			cmd += c.opHandler.Concat(ops.Cs(regRa), ops.Cs(regRb), ops.Cs(ops.RX))
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
	return c.opHandler.Move(ops.Cs(expr.Name.Lexeme), ops.Cs(ops.RX))
}

func (c *Compiler) VisitFieldAccess(expr expressions.FieldAccessExpr) string {
	cmd := "### BEGIN Struct field access operation ###\n"
	cmd += expr.Source.Accept(c)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(ops.RA))
	cmd += c.opHandler.StructGet(ops.Cs(ops.RA), expr.Field.Lexeme, ops.Cs(ops.RX))
	cmd += "### END   Struct field access operation ###\n"
	return cmd
}

func (c *Compiler) VisitFunctionCall(expr expressions.FunctionCallExpr) string {
	cmd := ""
	for i, arg := range expr.Arguments {
		cmd += arg.Accept(c)
		argName := c.functions[expr.Name.Lexeme].Args[i].Name
		cmd += c.opHandler.LoadArg(expr.Name.Lexeme, argName, ops.Cs(ops.RX))
	}
	if expr.ReturnType() != types.VoidType {
		cmd += c.opHandler.Call(expr.Name.Lexeme, ops.Cs(ops.RX))
	} else {
		cmd += c.opHandler.Call(expr.Name.Lexeme, "")
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
			cmd += c.opHandler.NegateNumber(ops.Cs(ops.RX))
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
	regRa := c.newRegister(ops.RA)
	cmd += c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(regRa))
	regRb := c.newRegister(ops.RB)
	cmd += c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(regRb))

	cmd += "### BEGIN Logical operation ###\n"
	leftSide += "###       Logical operation left side ###\n"
	leftSide += expr.Left.Accept(c)
	leftSide += c.opHandler.Move(ops.Cs(ops.RX), regRa)
	rightSide += "###       Logical operation right side ###\n"
	rightSide += expr.Right.Accept(c)
	rightSide += c.opHandler.Move(ops.Cs(ops.RX), regRb)
	rightSide += c.opHandler.MoveScore(regRb, regRb)

	cmd += leftSide
	switch expr.Operator.Type {
	case tokens.And:
		// If left side is false, return false
		evalRightSide := ""
		cmd += c.opHandler.MoveScore(regRa, regRa)
		cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), true, c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(ops.RX)))
		evalRightSide += rightSide
		evalRightSide += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), true, c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(ops.RX)))
		evalRightSide += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), false, c.opHandler.Move(regRb, ops.Cs(ops.RX)))
		cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), false, evalRightSide)
	case tokens.Or:
		// If left side is true, return true
		evalRightSide := ""
		cmd += c.opHandler.MoveScore(regRa, regRa)
		cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), false, c.opHandler.Move(regRa, ops.Cs(ops.RX)))
		evalRightSide += rightSide
		evalRightSide += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), true, c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(ops.RX)))
		evalRightSide += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRb, c.Namespace), false, c.opHandler.Move(regRb, ops.Cs(ops.RX)))
		cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 0", regRa, c.Namespace), true, evalRightSide)
	default:
		c.error(expr.SourceLocation, "Invalid operator for logical expressions")
	}

	cmd += "### END   Logical operation ###\n"
	return cmd
}

func (c *Compiler) VisitSlice(expr expressions.SliceExpr) string {
	regIndexStart := c.newRegister(ops.RA)
	regIndexEnd := c.newRegister(ops.RB)

	cmd := ""
	cmd += c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(regIndexStart))
	cmd += c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(regIndexEnd))
	cmd += "### BEGIN String slice operation ###\n"
	cmd += "###       Accept start index ###\n"
	cmd += expr.StartIndex.Accept(c)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regIndexStart))

	if expr.EndIndex == nil {
		cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regIndexEnd))
		cmd += c.opHandler.Inc(ops.Cs(regIndexEnd))
	} else {
		cmd += "###       Accept end index ###\n"
		cmd += expr.EndIndex.Accept(c)
		cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(regIndexEnd))
	}
	cmd += expr.TargetExpr.Accept(c)

	// Check index bounds

	cmd += "###       Cheking index bounds ###\n"
	lenReg := c.newRegister(ops.RX)
	cmd += c.opHandler.MakeConst(nbt.NewInt(0), ops.Cs(lenReg))

	cmd += c.opHandler.Size(ops.Cs(ops.RX), ops.Cs(lenReg))

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
	switch expr.TargetExpr.ReturnType().(type) {
	case types.PrimitiveTypeStruct:
		switch expr.TargetExpr.ReturnType() {
		case types.StringType:
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
			cmd += "###       Slice string ###\n"
			cmd += c.opHandler.SliceString(ops.Cs(ops.RX), ops.Cs(regIndexStart), ops.Cs(regIndexEnd), ops.Cs(ops.RX))
		}
	case types.ListTypeStruct:
		cmd += c.opHandler.ExecCond(
			fmt.Sprintf("score %s %s >= %s %s", ops.Cs(regIndexStart), c.Namespace, ops.Cs(lenReg), c.Namespace),
			true,
			c.opHandler.Exception("Index out of bounds"),
		)
		cmd += "###       Index list ###\n"
		if expr.EndIndex != nil {
			c.error(expr.SourceLocation, "List slicing is not supported")
			return ""
		}
		cmd += c.opHandler.MakeIndex(ops.Cs(regIndexStart), ops.Cs(regIndexStart))
		cmd += c.opHandler.PathGet(ops.Cs(ops.RX), ops.Cs(regIndexStart), ops.Cs(ops.RX))
	}
	cmd += "### END   String slice operation ###\n"
	return cmd
}

func (c *Compiler) VisitList(expr expressions.ListExpr) string {
	cmd := "### BEGIN List init operation ###\n"
	regList := ops.Cs(c.newRegister(ops.RX))
	cmd += c.opHandler.MakeConst(nbt.NewInt(0), regList)
	cmd += c.opHandler.MakeConst(nbt.NewList(), regList)
	for _, elem := range expr.Elements {
		cmd += elem.Accept(c)
		cmd += c.opHandler.AppendList(regList, ops.Cs(ops.RX))
	}
	cmd += c.opHandler.Move(regList, ops.Cs(ops.RX))
	cmd += "### END   List operation ###\n"
	return cmd
}

func (c *Compiler) VisitStruct(expr expressions.StructExpr) string {
	cmd := "### BEGIN Struct init operation ###\n"
	regStruct := ops.Cs(c.newRegister(ops.RX))
	cmd += c.opHandler.MakeConst(expr.StructType.ToNBT(), regStruct)
	structFields := expr.StructType.GetFieldNames()
	for i, arg := range expr.Args {
		cmd += arg.Accept(c)
		cmd += c.opHandler.StructSet(ops.Cs(ops.RX), structFields[i], regStruct)
	}
	cmd += c.opHandler.Move(regStruct, ops.Cs(ops.RX))
	cmd += "### END   Struct operation ###\n"
	return cmd
}
