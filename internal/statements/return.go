package statements

import "github.com/Kolterdyx/mcbasic/internal/expressions"

type ReturnStmt struct {
	Stmt

	Expression expressions.Expr
}

func (s ReturnStmt) Accept(v StmtVisitor) string {
	return v.VisitReturn(s)
}

func (s ReturnStmt) StmtType() StmtType {
	return ReturnStmtType
}
