package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type FieldAccessExpr struct {
	Expr
	interfaces.SourceLocation

	Source Expr
	Field  tokens.Token
	Type   types.ValueType
}

func (v FieldAccessExpr) Accept(visitor ExprVisitor) interfaces.IRCode {
	return visitor.VisitFieldAccess(v)
}

func (v FieldAccessExpr) ExprType() ExprType {
	return FieldAccessExprType
}

func (v FieldAccessExpr) ReturnType() types.ValueType {
	return v.Type
}

func (v FieldAccessExpr) ToString() string {
	return v.Source.ToString() + "." + v.Field.Lexeme
}
