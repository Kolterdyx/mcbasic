package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ListExpr struct {
	Expr
	interfaces.SourceLocation

	Elements  []Expr
	ValueType interfaces.ValueType
}

func (l ListExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitList(l)
}

func (l ListExpr) TType() ExprType {
	return ListExprType
}

func (l ListExpr) ReturnType() interfaces.ValueType {
	return l.ValueType
}
