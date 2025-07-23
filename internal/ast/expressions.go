package ast

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"reflect"
)

type ExpressionVisitor interface {
	VisitBinary(b *BinaryExpr) any
	VisitGrouping(g *GroupingExpr) any
	VisitLiteral(l *LiteralExpr) any
	VisitUnary(u *UnaryExpr) any
	VisitVariable(v *VariableExpr) any
	VisitDotAccess(v *DotAccessExpr) any
	VisitCall(f *CallExpr) any
	VisitLogical(l *LogicalExpr) any
	VisitSlice(s *SliceExpr) any
	VisitList(s *ListExpr) any
	VisitFunctionDeclaration(f FunctionDeclarationExpr) any
}

type Expr interface {
	Node
	Accept(v ExpressionVisitor) any
	ToString() string
	SetResolvedType(typ types.ValueType)
	GetResolvedType() types.ValueType
}

func AcceptExpr[T any](expr Expr, v ExpressionVisitor) T {
	if expr == nil {
		var zero T
		return zero
	}
	res := expr.Accept(v)
	if res == nil {
		var zero T
		return zero
	}
	val, ok := res.(T)
	if !ok {
		panic(fmt.Sprintf("unexpected type: got %v, want %v", reflect.TypeOf(res), reflect.TypeFor[T]()))
	}
	return val
}

type ResolvedType struct {
	typ types.ValueType
}

func (r *ResolvedType) SetResolvedType(typ types.ValueType) {
	r.typ = typ
}

func (r *ResolvedType) GetResolvedType() types.ValueType {
	return r.typ
}
