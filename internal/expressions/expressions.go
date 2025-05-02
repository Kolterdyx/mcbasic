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
	StructExprType       ExprType = "Struct"
)

type ExprVisitor interface {
	VisitBinary(b BinaryExpr) (cmd string)
	VisitGrouping(g GroupingExpr) (cmd string)
	VisitLiteral(l LiteralExpr) (cmd string)
	VisitUnary(u UnaryExpr) (cmd string)
	VisitVariable(v VariableExpr) (cmd string)
	VisitFieldAccess(v FieldAccessExpr) (cmd string)
	VisitFunctionCall(f FunctionCallExpr) (cmd string)
	VisitLogical(l LogicalExpr) (cmd string)
	VisitSlice(s SliceExpr) (cmd string)
	VisitList(s ListExpr) (cmd string)
	VisitStruct(s StructExpr) (cmd string)
}

type Expr interface {
	Accept(v ExprVisitor) (cmd string)
	ExprType() ExprType
	ReturnType() types.ValueType
}
