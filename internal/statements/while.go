package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type WhileStmt struct {
	Stmt

	Condition expressions.Expr
	Body      BlockStmt
}

func (w WhileStmt) Accept(v StmtVisitor) interfaces.IRCode {
	return v.VisitWhile(w)
}

func (w WhileStmt) Type() ast.NodeType {
	return ast.WhileStatement
}

func (w WhileStmt) ToString() string {
	body := ""
	for _, stmt := range w.Body.Statements {
		body += stmt.ToString() + "\n"
	}
	return "while " + w.Condition.ToString() + " {\n" + body + "}"
}
