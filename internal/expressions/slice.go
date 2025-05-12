package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type SliceExpr struct {
	StartIndex Expr
	EndIndex   Expr
	TargetExpr Expr

	interfaces.SourceLocation
	Expr
}

func (s SliceExpr) Accept(v ExprVisitor) interfaces.IRCode {
	return v.VisitSlice(s)
}

func (s SliceExpr) Type() ast.NodeType {
	return ast.SliceExpression
}

func (s SliceExpr) ToString() string {
	if s.StartIndex == nil && s.EndIndex == nil {
		return s.TargetExpr.ToString()
	}

	result := s.TargetExpr.ToString() + "["
	if s.StartIndex != nil {
		result += s.StartIndex.ToString()
	}
	if s.EndIndex != nil {
		result += ":" + s.EndIndex.ToString()
	}
	result += "]"
	return result
}
