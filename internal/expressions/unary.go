package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type UnaryExpr struct {
	Expr
	interfaces.SourceLocation

	Operator   tokens.Token
	Expression Expr
}

func (u UnaryExpr) Accept(v ExprVisitor) string {
	return v.VisitUnary(u)
}

func (u UnaryExpr) ExprType() ExprType {
	return UnaryExprType
}

func (u UnaryExpr) ReturnType() interfaces.ValueType {
	return types.IntType
}
