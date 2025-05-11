package statements

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type BlockStmt struct {
	Stmt

	Statements []Stmt
}

func (b BlockStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitBlock(b)
}

func (b BlockStmt) StmtType() StmtType {
	return BlockStmtType
}
