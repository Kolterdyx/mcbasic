package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	"strconv"
	"strings"
)

func (c *Compiler) VisitFunctionDeclaration(stmt statements.FunctionDeclarationStmt) interface{} {
	c.currentFunction = stmt
	c.currentScope = stmt.Name.Lexeme
	source := c.opHandler.RegLoad(strconv.Itoa(internal.FixedPointMagnitude), ops.RCF)
	c.opHandler.Scope = c.currentScope

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

func (c *Compiler) VisitExec(stmt statements.ExecStmt) interface{} {
	return c.opHandler.Exec(stmt.Command)
}

func (c *Compiler) VisitIf(stmt statements.IfStmt) interface{} {
	ifReg := c.newRegister(ops.RX)
	cmd := ""
	cmd += stmt.Condition.Accept(c).(string)
	cmd += c.opHandler.RegShift(ops.RX, ifReg)
	cmd += c.opHandler.And(ifReg, ifReg)
	thenBranch := stmt.ThenBranch.Accept(c).(string)
	elseBranch := stmt.ElseBranch.Accept(c).(string)
	cond := fmt.Sprintf("score %s %s matches 1..", ifReg, c.Namespace)
	cmd += c.opHandler.Ift(cond, strings.Split(thenBranch, "\n"))
	cmd += c.opHandler.Ifn(cond, strings.Split(elseBranch, "\n"))
	return cmd
}

func (c *Compiler) VisitWhile(stmt statements.WhileStmt) interface{} {
	condReg := c.newRegister(ops.RX)
	cmd := ""
	cmd += stmt.Condition.Accept(c).(string)
	cmd += c.opHandler.RegShift(ops.RX, condReg)
	cmd += c.opHandler.And(condReg, condReg)
	loop := stmt.Body.Accept(c).(string)
	cond := fmt.Sprintf("score %s %s matches 1..", condReg, c.Namespace)
	cmd += c.opHandler.Ift(cond, strings.Split(loop, "\n"))
	return cmd
}
