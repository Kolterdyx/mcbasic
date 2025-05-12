package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type VariableExpr struct {
	Expr
	interfaces.SourceLocation

	Name tokens.Token
}

func (v VariableExpr) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitVariable(v)
}

func (v VariableExpr) Type() NodeType {
	return VariableExpression
}

func (v VariableExpr) ToString() string {
	return v.Name.Lexeme
}
