package expressions

import "github.com/Kolterdyx/mcbasic/internal/tokens"

type VariableExpr struct {
	Expr

	Name tokens.Token
}

func (v VariableExpr) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariable(v)
}

func (v VariableExpr) Type() ExprType {
	return VariableExprType
}
