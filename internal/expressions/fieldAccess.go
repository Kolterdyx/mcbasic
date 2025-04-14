package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type FieldAccessExpr struct {
	Expr
	interfaces.SourceLocation

	Source Expr
	Field  tokens.Token
	Type   interfaces.ValueType
}

func (v FieldAccessExpr) Accept(visitor ExprVisitor) string {
	return visitor.VisitFieldAccess(v)
}

func (v FieldAccessExpr) ExprType() ExprType {
	return FieldAccessExprType
}

func (v FieldAccessExpr) ReturnType() interfaces.ValueType {
	return v.Type
}
