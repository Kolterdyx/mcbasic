package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type SliceExpr struct {
	StartIndex Expr
	EndIndex   Expr
	TargetExpr Expr

	interfaces.SourceLocation
	Expr
}

func (s SliceExpr) Accept(v ExprVisitor) string {
	return v.VisitSlice(s)
}

func (s SliceExpr) ExprType() ExprType {
	return SliceExprType
}

func (s SliceExpr) ReturnType() interfaces.ValueType {
	if s.StartIndex != nil && s.EndIndex == nil {
		return getReturnIndexType(s.TargetExpr.ReturnType())
	}
	return s.TargetExpr.ReturnType()
}

func getReturnIndexType(valueType interfaces.ValueType) interfaces.ValueType {
	switch reflect.TypeOf(valueType) {
	case reflect.TypeOf(types.ListTypeStruct{}):
		return valueType.Primitive()
	default:
		log.Errorf("Can't index type: %v", valueType)
		return types.ErrorType
	}
}
