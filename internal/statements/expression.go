package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type ExpressionStmt struct {
	Stmt

	Expression expressions.Expr
}

func (e ExpressionStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitExpression(e)
}

func (e ExpressionStmt) StmtType() StmtType {
	return ExpressionStmtType
}
