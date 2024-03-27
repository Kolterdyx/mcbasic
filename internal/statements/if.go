package statements

import "github.com/Kolterdyx/mcbasic/internal/expressions"

type IfStmt struct {
	Stmt

	Condition  expressions.Expr
	ThenBranch BlockStmt
	ElseBranch BlockStmt
}

func (i IfStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitIf(i)
}

func (i IfStmt) TType() StmtType {
	return IfStmtType
}
