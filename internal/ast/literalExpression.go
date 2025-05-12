package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type LiteralExpr struct {
	Expr
	interfaces.SourceLocation

	Value     nbt.Value
	ValueType types.ValueType
}

func (l LiteralExpr) Accept(v ExpressionVisitor) any {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) Type() NodeType {
	return LiteralExpression
}

func (l LiteralExpr) ReturnType() types.ValueType {
	return l.ValueType
}

func (l LiteralExpr) ToString() string {
	return l.Value.ToString()
}

func (l LiteralExpr) GetSourceLocation() interfaces.SourceLocation {
	return l.SourceLocation
}
