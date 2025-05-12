package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
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

func (g GroupingExpr) ToString() string {
	return "(" + g.Expression.ToString() + ")"
}
