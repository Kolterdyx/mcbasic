package ast

import (
	"fmt"
	"reflect"
)

type ExpressionVisitor interface {
	VisitBinary(b BinaryExpr) any
	VisitGrouping(g GroupingExpr) any
	VisitLiteral(l LiteralExpr) any
	VisitUnary(u UnaryExpr) any
	VisitVariable(v VariableExpr) any
	VisitFieldAccess(v FieldAccessExpr) any
	VisitFunctionCall(f FunctionCallExpr) any
	VisitLogical(l LogicalExpr) any
	VisitSlice(s SliceExpr) any
	VisitList(s ListExpr) any
}

type Expr interface {
	Node
	Accept(v ExpressionVisitor) any
	ToString() string
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
