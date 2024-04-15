package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type UnaryExpr struct {
	Expr
	interfaces.SourceLocation

	Operator   tokens.Token
	Expression Expr
}

func (u UnaryExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitUnary(u)
}

func (u UnaryExpr) TType() ExprType {
	return UnaryExprType
}

func (u UnaryExpr) ReturnType() ValueType {
	return NumberType
}
