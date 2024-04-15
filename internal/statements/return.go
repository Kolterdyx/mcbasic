package statements

import "github.com/Kolterdyx/mcbasic/internal/expressions"

type ReturnStmt struct {
	Stmt

	Expression expressions.Expr
}

func (s ReturnStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitReturn(s)
}

func (s ReturnStmt) TType() StmtType {
	return ReturnStmtType
}
