package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type ListExpr struct {
	Expr
	interfaces.SourceLocation

	Elements  []Expr
	ValueType types.ValueType
}

func (l ListExpr) Accept(v ExpressionVisitor) any {
	return v.VisitList(l)
}

func (l ListExpr) Type() NodeType {
	return ListExpression
}

func (l ListExpr) ToString() string {
	if len(l.Elements) == 0 {
		return "[]"
	}

	result := "["
	for i, element := range l.Elements {
		if i > 0 {
			result += ", "
		}
		result += element.ToString()
	}
	result += "]"
	return result
}

func (l ListExpr) GetSourceLocation() interfaces.SourceLocation {
	return l.SourceLocation
}
