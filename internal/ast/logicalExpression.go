package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type LogicalExpr struct {
	Left     Expr
	Operator tokens.Token
	Right    Expr

	Expr
}

func (l LogicalExpr) Accept(v ExpressionVisitor) any {
	return v.VisitLogical(l)
}

func (l LogicalExpr) Type() NodeType {
	return LogicalExpression
}

func (l LogicalExpr) ToString() string {
	return "(" + l.Left.ToString() + " " + l.Operator.Lexeme + " " + l.Right.ToString() + ")"
}

func (l LogicalExpr) GetSourceLocation() interfaces.SourceLocation {
	return l.GetSourceLocation()
}
