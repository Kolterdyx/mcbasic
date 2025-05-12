package ast

type ExpressionVisitor interface {
	VisitBinary(b BinaryExpr) any
	VisitGrouping(g GroupingExpr) any
	VisitLiteral(l LiteralExpr) any
	VisitUnary(u UnaryExpr) any
	VisitVariable(v VariableExpr) any
	VisitFieldAccess(v FieldAccessExpr) any
	VisitFunctionCall(f FunctionCallExpr) any
	VisitLogical(l LogicalExpr) any
	VisitSlice(s SliceExpr) any
	VisitList(s ListExpr) any
}

type Expr interface {
	Node
	Accept(v ExpressionVisitor) any
	ToString() string
	Validate() error
}

func AcceptExpr[T any](expr Expr, v ExpressionVisitor) T {
	return expr.Accept(v).(T)
}
