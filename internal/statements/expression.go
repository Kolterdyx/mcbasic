package statements

import "github.com/Kolterdyx/mcbasic/internal/expressions"

type ExpressionStmt struct {
	Stmt

	Expression expressions.Expr
}

func (e ExpressionStmt) Accept(visitor StmtVisitor) string {
	return visitor.VisitExpression(e)
}

func (e ExpressionStmt) TType() StmtType {
	return ExpressionStmtType
}
