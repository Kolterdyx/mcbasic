package expressions

import "github.com/Kolterdyx/mcbasic/internal/tokens"

type FunctionCallExpr struct {
	Expr

	Name      tokens.Token
	Arguments []Expr
}

func (f FunctionCallExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitFunctionCall(f)
}

func (f FunctionCallExpr) Type() ExprType {
	return FunctionCallExprType
}
