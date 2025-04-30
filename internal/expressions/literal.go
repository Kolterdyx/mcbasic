package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type LiteralExpr struct {
	Expr
	interfaces.SourceLocation

	Value     string
	ValueType interfaces.ValueType
}

func (l LiteralExpr) Accept(v ExprVisitor) string {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) ExprType() ExprType {
	return LiteralExprType
}

func (l LiteralExpr) ReturnType() interfaces.ValueType {
	return l.ValueType
}
