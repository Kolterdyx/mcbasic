package ast

type BlockStmt struct {
	Statement

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
