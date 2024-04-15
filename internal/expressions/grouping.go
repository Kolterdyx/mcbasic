package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type GroupingExpr struct {
	Expr
	interfaces.SourceLocation

	Expression Expr
}

func (g GroupingExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitGrouping(g)
}

func (g GroupingExpr) TType() ExprType {
	return GroupingExprType
}
