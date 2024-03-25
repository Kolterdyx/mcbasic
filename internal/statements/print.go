package statements

import "github.com/Kolterdyx/mcbasic/internal/expressions"

type PrintStmt struct {
	Stmt

	Expression expressions.Expr
}

func (p PrintStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrint(p)
}

func (p PrintStmt) Type() StmtType {
	return PrintStmtType
}
