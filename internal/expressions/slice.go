package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
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

func (s SliceExpr) ReturnType() types.ValueType {
	if s.StartIndex != nil && s.EndIndex == nil {
		return getReturnIndexType(s.TargetExpr.ReturnType())
	}
	return s.TargetExpr.ReturnType()
}

func getReturnIndexType(valueType types.ValueType) types.ValueType {
	switch valueType.(type) {
	case types.ListTypeStruct:
		return valueType.(types.ListTypeStruct).ContentType
	default:
		log.Errorf("Can't index type: %v", valueType)
		return nil
	}
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
