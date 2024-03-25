package expressions

type ValueType string

const (
	NumberType ValueType = "number"
	StringType ValueType = "string"
)

type LiteralExpr struct {
	Expr

	Value     interface{}
	ValueType ValueType
}

func (l LiteralExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) Type() ExprType {
	return LiteralExprType
}
