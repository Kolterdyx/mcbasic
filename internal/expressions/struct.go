package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type StructExpr struct {
	Expr
	interfaces.SourceLocation

	Args       []Expr
	StructType types.StructTypeStruct
}

func (s StructExpr) Accept(v ExprVisitor) interfaces.IRCode {
	return v.VisitStruct(s)
}

func (s StructExpr) ExprType() ExprType {
	return StructExprType
}

func (s StructExpr) ReturnType() types.ValueType {
	return s.StructType
}
