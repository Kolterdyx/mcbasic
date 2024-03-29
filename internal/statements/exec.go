package statements

type ExecStmt struct {
	Stmt

	Command string
}

func (e ExecStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExec(e)
}

func (e ExecStmt) TType() StmtType {
	return ExecStmtType
}
