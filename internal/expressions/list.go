package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type ListExpr struct {
	Expr
	interfaces.SourceLocation

	Elements  []Expr
	ValueType types.ValueType
}

func (l ListExpr) Accept(v ExprVisitor) string {
	return v.VisitList(l)
}

func (l ListExpr) ExprType() ExprType {
	return ListExprType
}

func (l ListExpr) ReturnType() types.ValueType {
	return l.ValueType
}
