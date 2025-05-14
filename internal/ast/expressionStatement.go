package ast

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ExpressionStmt struct {
	Expression Expr
}

func (e ExpressionStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitExpression(e)
}

func (e ExpressionStmt) Type() NodeType {
	return ExpressionStatement
}

func (e ExpressionStmt) ToString() string {
	return e.Expression.ToString() + ";"
}

func (e ExpressionStmt) GetSourceLocation() interfaces.SourceLocation {
	return e.Expression.GetSourceLocation()
}
