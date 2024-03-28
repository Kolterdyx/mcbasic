package compiler

import (
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	"strconv"
)

func (c *Compiler) VisitFunctionDeclaration(stmt statements.FunctionDeclarationStmt) interface{} {
	c.currentFunction = stmt
	source := c.opHandler.RegLoad(strconv.Itoa(internal.FixedPointMagnitude), ops.RCF)

	// Store arguments in variables
	for i, arg := range stmt.Parameters {
		source += c.opHandler.SetMacro(arg.Lexeme, c.opHandler.Macro(arg.Lexeme))
		source += c.opHandler.RegLoad(strconv.Itoa(i), ops.RA)
		source += c.opHandler.RegShift(ops.RA, ops.RX)
	}

	source += stmt.Body.Accept(c).(string)
	source = c.opHandler.MacroReplace(source)
	c.createFunction(stmt.Name.Lexeme, source)
	return ""
}

func (c *Compiler) VisitVariableDeclaration(stmt statements.VariableDeclarationStmt) interface{} {
	cmd := ""
	if stmt.Initializer != nil {
		cmd += stmt.Initializer.Accept(c).(string)
		cmd += c.opHandler.RegSave(stmt.Name.Lexeme, ops.RX)
		return cmd
	}
	return c.opHandler.Set(stmt.Name.Lexeme, "0")
}

func (c *Compiler) VisitVariableAssignment(stmt statements.VariableAssignmentStmt) interface{} {
	cmd := stmt.Value.Accept(c).(string)
	cmd += c.opHandler.Set(stmt.Name.Lexeme, ops.RX)
	return cmd
}

func (c *Compiler) VisitExpression(stmt statements.ExpressionStmt) interface{} {
	return stmt.Expression.Accept(c)
}

func (c *Compiler) VisitBlock(stmt statements.BlockStmt) interface{} {
	cmd := ""
	for _, s := range stmt.Statements {
		cmd += s.Accept(c).(string)
	}
	return cmd
}

func (c *Compiler) VisitPrint(stmt statements.PrintStmt) interface{} {
	cmd := stmt.Expression.Accept(c).(string)
	if stmt.Expression.Type() == expressions.LiteralExprType {
		cmd += c.opHandler.ArgLoad("builtin/print", "text", ops.RX)
	} else {
		cmd += c.opHandler.RegSave(ops.RX, ops.RX)
		cmd += c.opHandler.ArgLoad("builtin/print", "text", ops.RX)
	}
	cmd += c.opHandler.Call("builtin/print")
	return cmd
}
