package compiler

import (
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
)

func (c *Compiler) VisitFunctionDeclaration(stmt statements.FunctionDeclarationStmt) interface{} {
	c.currentFunction = stmt
	c.currentScope = stmt.Name.Lexeme
	c.scope[stmt.Name.Lexeme] = []string{}

	var source = stmt.Body.Accept(c).(string)

	args := make([]statements.FuncArg, 0)
	for _, arg := range stmt.Parameters {
		args = append(args, statements.FuncArg{Name: arg.Name, Type: arg.Type})
		c.scope[stmt.Name.Lexeme] = append(c.scope[stmt.Name.Lexeme], arg.Name)
	}

	c.createFunction(stmt.Name.Lexeme, source, args)
	return ""
}

func (c *Compiler) VisitBlock(stmt statements.BlockStmt) interface{} {
	cmd := ""
	for _, s := range stmt.Statements {
		cmd += s.Accept(c).(string)
	}
	return cmd
}

func (c *Compiler) VisitExpression(stmt statements.ExpressionStmt) interface{} {
	return stmt.Expression.Accept(c)
}

func (c *Compiler) VisitReturn(stmt statements.ReturnStmt) interface{} {
	cmd := ""
	cmd += stmt.Expression.Accept(c).(string)
	cmd += c.opHandler.Move(ops.RX, ops.RET)
	return cmd
}

func (c *Compiler) VisitVariableDeclaration(stmt statements.VariableDeclarationStmt) interface{} {
	cmd := ""
	if stmt.Initializer != nil {
		cmd += stmt.Initializer.Accept(c).(string)
	}
	cmd += c.opHandler.Move(ops.RX, stmt.Name.Lexeme)
	c.scope[c.currentScope] = append(c.scope[c.currentScope], stmt.Name.Lexeme)
	return cmd
}
