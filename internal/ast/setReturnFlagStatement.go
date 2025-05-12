package ast

type SetReturnFlagStmt struct {
	Statement
}

func (s SetReturnFlagStmt) Accept(v StatementVisitor) any {
	return v.VisitSetReturnFlag(s)
}

func (s SetReturnFlagStmt) Type() NodeType {
	return SetReturnFlagStatement
}

func (s SetReturnFlagStmt) ToString() string {
	return "retf;"
}
