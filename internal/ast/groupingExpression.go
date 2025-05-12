package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type GroupingExpr struct {
	Expr
	interfaces.SourceLocation

	Expression Expr
}

func (g GroupingExpr) Accept(v ExpressionVisitor) any {
	return v.VisitGrouping(g)
}

func (g GroupingExpr) Type() NodeType {
	return GroupingExpression
}

func (g GroupingExpr) ToString() string {
	return "(" + g.Expression.ToString() + ")"
}
