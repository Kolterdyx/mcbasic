package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type IfStmt struct {
	Stmt

	Condition  expressions.Expr
	ThenBranch BlockStmt
	ElseBranch *BlockStmt
}

func (i IfStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitIf(i)
}

func (i IfStmt) StmtType() StmtType {
	return IfStmtType
}
