package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type FunctionCallExpr struct {
	Expr
	interfaces.SourceLocation

	Name      tokens.Token
	Arguments []Expr
	Type      types.ValueType
}

func (f FunctionCallExpr) Accept(visitor ExprVisitor) string {
	return visitor.VisitFunctionCall(f)
}

func (f FunctionCallExpr) ExprType() ExprType {
	return FunctionCallExprType
}

func (f FunctionCallExpr) ReturnType() types.ValueType {
	return f.Type
}
