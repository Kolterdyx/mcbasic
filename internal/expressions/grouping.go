package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type GroupingExpr struct {
	Expr
	interfaces.SourceLocation

	Expression Expr
}

func (g GroupingExpr) Accept(v ExprVisitor) interfaces.IRCode {
	return v.VisitGrouping(g)
}

func (g GroupingExpr) ExprType() ExprType {
	return GroupingExprType
}

func (g GroupingExpr) ReturnType() types.ValueType {
	return g.Expression.ReturnType()
}
