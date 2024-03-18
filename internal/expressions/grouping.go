package expressions

type GroupingExpr struct {
	Expr

	Expression Expr
}

func (g GroupingExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitGrouping(g)
}

func (g GroupingExpr) Type() ExprType {
	return GroupingExprType
}
