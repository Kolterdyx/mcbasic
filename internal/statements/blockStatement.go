package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type BlockStmt struct {
	Stmt

	Statements []Stmt
}

func (b BlockStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitBlock(b)
}

func (b BlockStmt) Type() ast.NodeType {
	return ast.BlockStatement
}

func (b BlockStmt) ToString() string {
	body := ""
	for _, stmt := range b.Statements {
		body += stmt.ToString() + "\n"
	}
	return body
}
