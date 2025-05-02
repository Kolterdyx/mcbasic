package statements

type ScoreStmt struct {
	Stmt
	Target string
	Score  int64
}

func (s ScoreStmt) Accept(visitor StmtVisitor) string {
	return visitor.VisitScore(s)
}

func (s ScoreStmt) StmtType() StmtType {
	return ScoreStmtType
}
