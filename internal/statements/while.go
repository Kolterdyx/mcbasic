package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type WhileStmt struct {
	Stmt

	Condition expressions.Expr
	Body      BlockStmt
}

func (w WhileStmt) Accept(v StmtVisitor) interfaces.IRCode {
	return v.VisitWhile(w)
}

func (w WhileStmt) StmtType() StmtType {
	return WhileStmtType
}
