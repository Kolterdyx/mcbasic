package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type SliceExpr struct {
	StartIndex Expr
	EndIndex   Expr
	TargetExpr Expr

	interfaces.SourceLocation
	Expr
}

func (s SliceExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitSlice(s)
}

func (s SliceExpr) TType() ExprType {
	return SliceExprType
}

func (s SliceExpr) ReturnType() ValueType {
	return s.TargetExpr.ReturnType()
}
