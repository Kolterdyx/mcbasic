package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type VariableExpr struct {
	Expr
	interfaces.SourceLocation

	Name tokens.Token
	Type interfaces.ValueType
}

func (v VariableExpr) Accept(visitor ExprVisitor) string {
	return visitor.VisitVariable(v)
}

func (v VariableExpr) ExprType() ExprType {
	return VariableExprType
}

func (v VariableExpr) ReturnType() interfaces.ValueType {
	return v.Type
}
