package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
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

func (g GroupingExpr) Type() ast.NodeType {
	return ast.GroupingExpression
}

func (g GroupingExpr) ReturnType() types.ValueType {
	return g.Expression.ReturnType()
}

func (g GroupingExpr) ToString() string {
	return "(" + g.Expression.ToString() + ")"
}
