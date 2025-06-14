package ast

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type WhileStmt struct {
	Condition Expr
	Body      BlockStmt
}

func (w WhileStmt) Accept(v StatementVisitor) any {
	return v.VisitWhile(w)
}

func (w WhileStmt) Type() NodeType {
	return WhileStatement
}

func (w WhileStmt) ToString() string {
	body := ""
	for _, stmt := range w.Body.Statements {
		body += stmt.ToString() + "\n"
	}
	return "while " + w.Condition.ToString() + " {\n" + body + "}"
}

func (w WhileStmt) GetSourceLocation() interfaces.SourceLocation {
	return w.Condition.GetSourceLocation()
}
