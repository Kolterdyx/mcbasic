package ast

type IfStmt struct {
	Statement

	Condition  Expr
	ThenBranch BlockStmt
	ElseBranch *BlockStmt
}

func (i IfStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitIf(i)
}

func (i IfStmt) Type() NodeType {
	return IfStatement
}

func (i IfStmt) ToString() string {
	body := ""
	for _, stmt := range i.ThenBranch.Statements {
		body += stmt.ToString() + "\n"
	}
	if i.ElseBranch != nil {
		for _, stmt := range i.ElseBranch.Statements {
			body += stmt.ToString() + "\n"
		}
	}
	return "if " + i.Condition.ToString() + " {\n" + body + "}"
}
