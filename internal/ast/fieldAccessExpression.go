package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type FieldAccessExpr struct {
	Expr
	interfaces.SourceLocation

	Source Expr
	Field  tokens.Token
}

func (v FieldAccessExpr) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitFieldAccess(v)
}

func (v FieldAccessExpr) Type() NodeType {
	return FieldAccessExpression
}

func (v FieldAccessExpr) ToString() string {
	return v.Source.ToString() + "." + v.Field.Lexeme
}

func (v FieldAccessExpr) GetSourceLocation() interfaces.SourceLocation {
	return v.SourceLocation
}
