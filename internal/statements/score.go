package statements

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ScoreStmt struct {
	Stmt
	Target string
	Score  int64
}

func (s ScoreStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitScore(s)
}

func (s ScoreStmt) StmtType() StmtType {
	return ScoreStmtType
}
