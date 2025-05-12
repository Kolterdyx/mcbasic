package ast

type ExpressionStmt struct {
	Statement

	Expression Expr
}

func (e ExpressionStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitExpression(e)
}

func (e ExpressionStmt) Type() NodeType {
	return ExpressionStatement
}

func (e ExpressionStmt) ToString() string {
	return e.Expression.ToString() + ";"
}
