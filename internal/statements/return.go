package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type ReturnStmt struct {
	Stmt

	Expression expressions.Expr
}

func (s ReturnStmt) Accept(v StmtVisitor) interfaces.IRCode {
	return v.VisitReturn(s)
}

func (s ReturnStmt) StmtType() StmtType {
	return ReturnStmtType
}

func (s ReturnStmt) ToString() string {
	if s.Expression != nil {
		return "return " + s.Expression.ToString() + ";"
	}
	return "return;"
}
