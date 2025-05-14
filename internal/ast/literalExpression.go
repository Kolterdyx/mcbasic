package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type LiteralExpr struct {
	interfaces.SourceLocation

	Value nbt.Value
	ResolvedType
}

func NewLiteralExpr(value nbt.Value, typ types.ValueType, source interfaces.SourceLocation) *LiteralExpr {
	literal := &LiteralExpr{
		Value:          value,
		SourceLocation: source,
	}
	literal.SetResolvedType(typ)
	return literal
}

func (l *LiteralExpr) Accept(v ExpressionVisitor) any {
	return v.VisitLiteral(l)
}

func (l *LiteralExpr) Type() NodeType {
	return LiteralExpression
}

func (l *LiteralExpr) ToString() string {
	return l.Value.ToString()
}

func (l *LiteralExpr) GetSourceLocation() interfaces.SourceLocation {
	return l.SourceLocation
}
