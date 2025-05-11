package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type LiteralExpr struct {
	Expr
	interfaces.SourceLocation

	Value     string
	ValueType types.ValueType
}

func (l LiteralExpr) Accept(v ExprVisitor) interfaces.IRCode {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) ExprType() ExprType {
	return LiteralExprType
}

func (l LiteralExpr) ReturnType() types.ValueType {
	return l.ValueType
}
