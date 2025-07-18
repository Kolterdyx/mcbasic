package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type VariableExpr struct {
	Name tokens.Token
	ResolvedType
}

func (v *VariableExpr) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitVariable(v)
}

func (v *VariableExpr) Type() NodeType {
	return VariableExpression
}

func (v *VariableExpr) ToString() string {
	return v.Name.Lexeme
}

func (v *VariableExpr) GetSourceLocation() interfaces.SourceLocation {
	return v.Name.SourceLocation
}
