package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type GroupingExpr struct {
	Expr
	interfaces.SourceLocation

	Expression Expr
}

func (g GroupingExpr) Accept(v ExprVisitor) string {
	return v.VisitGrouping(g)
}

func (g GroupingExpr) ExprType() ExprType {
	return GroupingExprType
}

func (g GroupingExpr) ReturnType() interfaces.ValueType {
	return g.Expression.ReturnType()
}
