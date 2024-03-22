package statements

import "github.com/Kolterdyx/mcbasic/internal/expressions"

type WhileStmt struct {
	Stmt

	Condition expressions.Expr
	Body      BlockStmt
}

func (w WhileStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitWhile(w)
}

func (w WhileStmt) Type() StmtType {
	return WhileStmtType
}
