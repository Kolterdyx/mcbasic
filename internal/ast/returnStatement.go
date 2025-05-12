package ast

type ReturnStmt struct {
	Statement

	Expression Expr
}

func (s ReturnStmt) Accept(v StatementVisitor) any {
	return v.VisitReturn(s)
}

func (s ReturnStmt) Type() NodeType {
	return ReturnStatement
}

func (s ReturnStmt) ToString() string {
	if s.Expression != nil {
		return "return " + s.Expression.ToString() + ";"
	}
	return "return;"
}
