package statements

import "github.com/Kolterdyx/mcbasic/internal/expressions"

type ExecStmt struct {
	Stmt

	Expression expressions.LiteralExpr
}

func (e ExecStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExec(e)
}

func (e ExecStmt) TType() StmtType {
	return ExecStmtType
}
