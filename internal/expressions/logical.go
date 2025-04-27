package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type LogicalExpr struct {
	Left     Expr
	Operator tokens.Token
	Right    Expr

	interfaces.SourceLocation
	Expr
}

func (l LogicalExpr) Accept(v ExprVisitor) string {
	return v.VisitLogical(l)
}

func (l LogicalExpr) ExprType() ExprType {
	return LogicalExprType
}

func (l LogicalExpr) ReturnType() interfaces.ValueType {
	return types.IntType
}
