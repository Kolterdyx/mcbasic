package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ListExpr struct {
	Expr
	interfaces.SourceLocation

	Elements  []Expr
	ValueType interfaces.ValueType
}

func (l ListExpr) Accept(v ExprVisitor) string {
	return v.VisitList(l)
}

func (l ListExpr) ExprType() ExprType {
	return ListExprType
}

func (l ListExpr) ReturnType() interfaces.ValueType {
	return l.ValueType
}
