package statements

type BlockStmt struct {
	Stmt

	Statements []Stmt
}

func (b BlockStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitBlock(b)
}

func (b BlockStmt) TType() StmtType {
	return BlockStmtType
}
