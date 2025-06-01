package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type DotAccessExpr struct {
	interfaces.SourceLocation

	Source Expr
	Name   tokens.Token
	ResolvedType
}

func (v *DotAccessExpr) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitDotAccess(v)
}

func (v *DotAccessExpr) Type() NodeType {
	return FieldAccessExpression
}

func (v *DotAccessExpr) ToString() string {
	return v.Source.ToString() + "." + v.Name.Lexeme
}

func (v *DotAccessExpr) GetSourceLocation() interfaces.SourceLocation {
	return v.SourceLocation
}
