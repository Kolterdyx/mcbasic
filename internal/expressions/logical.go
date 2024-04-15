package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type LogicalExpr struct {
	Left     Expr
	Operator tokens.Token
	Right    Expr

	interfaces.SourceLocation
	Expr
}

func (l LogicalExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitLogical(l)
}

func (l LogicalExpr) TType() ExprType {
	return LogicalExprType
}

func (l LogicalExpr) ReturnType() ValueType {
	return NumberType
}
