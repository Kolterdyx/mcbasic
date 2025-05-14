package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type UnaryExpr struct {
	Operator   tokens.Token
	Expression Expr
}

func (u UnaryExpr) Accept(v ExpressionVisitor) any {
	return v.VisitUnary(u)
}

func (u UnaryExpr) Type() NodeType {
	return UnaryExpression
}

func (u UnaryExpr) ReturnType() types.ValueType {
	return types.IntType
}

func (u UnaryExpr) ToString() string {
	return "(" + u.Operator.Lexeme + " " + u.Expression.ToString() + ")"
}

func (u UnaryExpr) GetSourceLocation() interfaces.SourceLocation {
	return u.Operator.SourceLocation
}
