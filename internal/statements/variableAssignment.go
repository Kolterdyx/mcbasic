package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type VariableAssignmentStmt struct {
	Stmt

	Name  tokens.Token
	Index expressions.Expr
	Value expressions.Expr
}

func (v VariableAssignmentStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitVariableAssignment(v)
}

func (v VariableAssignmentStmt) TType() StmtType {
	return VariableAssignmentStmtType
}
