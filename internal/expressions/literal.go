package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ValueType string

const (
	NumberType ValueType = "number"
	StringType ValueType = "string"
)

type LiteralExpr struct {
	Expr
	interfaces.SourceLocation

	Value     interface{}
	ValueType ValueType
}

func (l LiteralExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) TType() ExprType {
	return LiteralExprType
}
