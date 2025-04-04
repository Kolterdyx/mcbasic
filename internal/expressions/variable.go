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

func (v VariableExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariable(v)
}

func (v VariableExpr) TType() ExprType {
	return VariableExprType
}

func (v VariableExpr) ReturnType() interfaces.ValueType {
	return v.Type
}
