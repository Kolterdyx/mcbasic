package ast

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type BlockStmt struct {
	Statement
	interfaces.SourceLocation

	Statements []Statement
}

func (b BlockStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitBlock(b)
}

func (b BlockStmt) Type() NodeType {
	return BlockStatement
}

func (b BlockStmt) ToString() string {
	body := ""
	for _, stmt := range b.Statements {
		body += stmt.ToString() + "\n"
	}
	return body
}

func (b BlockStmt) GetSourceLocation() interfaces.SourceLocation {
	return b.SourceLocation
}
