package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type VariableExpr struct {
	Expr
	interfaces.SourceLocation

	Name tokens.Token
	Type types.ValueType
}

func (v VariableExpr) Accept(visitor ExprVisitor) interfaces.IRCode {
	return visitor.VisitVariable(v)
}

func (v VariableExpr) ExprType() ExprType {
	return VariableExprType
}

func (v VariableExpr) ReturnType() types.ValueType {
	return v.Type
}
