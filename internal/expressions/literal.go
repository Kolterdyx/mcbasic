package expressions

type LiteralExpr struct {
	Expr

	Value interface{}
}

func (l LiteralExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) Type() ExprType {
	return LiteralExprType
}
