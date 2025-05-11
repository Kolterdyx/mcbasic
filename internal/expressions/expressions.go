package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type ExprType string

const (
	BinaryExprType       ExprType = "Binary"
	GroupingExprType     ExprType = "Grouping"
	LiteralExprType      ExprType = "Literal"
	UnaryExprType        ExprType = "Unary"
	FieldAccessExprType  ExprType = "FieldAccess"
	VariableExprType     ExprType = "Variable"
	FunctionCallExprType ExprType = "FunctionCall"
	LogicalExprType      ExprType = "Logical"
	SliceExprType        ExprType = "SliceString"
	ListExprType         ExprType = "List"
	StructExprType       ExprType = "Struct"
)

type ExprVisitor interface {
	VisitBinary(b BinaryExpr) interfaces.IRCode
	VisitGrouping(g GroupingExpr) interfaces.IRCode
	VisitLiteral(l LiteralExpr) interfaces.IRCode
	VisitUnary(u UnaryExpr) interfaces.IRCode
	VisitVariable(v VariableExpr) interfaces.IRCode
	VisitFieldAccess(v FieldAccessExpr) interfaces.IRCode
	VisitFunctionCall(f FunctionCallExpr) interfaces.IRCode
	VisitLogical(l LogicalExpr) interfaces.IRCode
	VisitSlice(s SliceExpr) interfaces.IRCode
	VisitList(s ListExpr) interfaces.IRCode
	VisitStruct(s StructExpr) interfaces.IRCode
}

type Expr interface {
	Accept(v ExprVisitor) interfaces.IRCode
	ExprType() ExprType
	ReturnType() types.ValueType
	ToString() string
}
