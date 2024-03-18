package expressions

import "github.com/Kolterdyx/mcbasic/internal/tokens"

type UnaryExpr struct {
	Expr

	Operator   tokens.Token
	Expression Expr
}

func (u UnaryExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitUnary(u)
}

func (u UnaryExpr) Type() ExprType {
	return UnaryExprType
}
