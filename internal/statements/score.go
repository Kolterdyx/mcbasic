package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"strconv"
)

type ScoreStmt struct {
	Stmt
	Target string
	Score  int64
}

func (s ScoreStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitScore(s)
}

func (s ScoreStmt) Type() ast.NodeType {
	return ast.ScoreStatement
}

func (s ScoreStmt) ToString() string {
	return "score " + s.Target + " = " + strconv.FormatInt(s.Score, 10) + ";"
}
