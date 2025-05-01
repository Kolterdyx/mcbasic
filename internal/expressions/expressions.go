package expressions

import (
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
)

type ExprVisitor interface {
	VisitBinary(b BinaryExpr) string
	VisitGrouping(g GroupingExpr) string
	VisitLiteral(l LiteralExpr) string
	VisitUnary(u UnaryExpr) string
	VisitVariable(v VariableExpr) string
	VisitFieldAccess(v FieldAccessExpr) string
	VisitFunctionCall(f FunctionCallExpr) string
	VisitLogical(l LogicalExpr) string
	VisitSlice(s SliceExpr) string
	VisitList(s ListExpr) string
}

type Expr interface {
	Accept(v ExprVisitor) string
	ExprType() ExprType
	ReturnType() types.ValueType
}
