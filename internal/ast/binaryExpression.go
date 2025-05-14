package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type BinaryExpr struct {
	Left     Expr
	Operator tokens.Token
	Right    Expr
}

func (b BinaryExpr) Accept(v ExpressionVisitor) any {
	return v.VisitBinary(b)
}

func (b BinaryExpr) Type() NodeType {
	return BinaryExpression
}

func (b BinaryExpr) ToString() string {
	return "(" + b.Left.ToString() + " " + b.Operator.Lexeme + " " + b.Right.ToString() + ")"
}

func (b BinaryExpr) GetSourceLocation() interfaces.SourceLocation {
	return b.Left.GetSourceLocation()
}
