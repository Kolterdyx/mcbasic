package compiler

import (
	"github.com/Kolterdyx/mcbasic/internal/statements"
)

func (c *Compiler) VisitFunctionDeclarationStmt(stmt statements.FunctionDeclarationStmt) interface{} {
	c.currentFunction = stmt
	return stmt.Body.Accept(c)
}

func (c *Compiler) VisitVariableDeclarationStmt(stmt statements.VariableDeclarationStmt) interface{} {
	if stmt.Initializer != nil {
		return c.opHandler.Set(stmt.Name.Lexeme, stmt.Initializer.Accept(c.ExprVisitor).(string))
	}
	return c.opHandler.Set(stmt.Name.Lexeme, "0")
}
