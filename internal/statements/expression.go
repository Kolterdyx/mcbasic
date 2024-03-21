package statements

import "github.com/Kolterdyx/mcbasic/internal/expressions"

type ExpressionStmt struct {
	Stmt

	Expression expressions.Expr
}

func (e ExpressionStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpression(e)
}

func (e ExpressionStmt) Type() StmtType {
	return ExpressionStmtType
}
