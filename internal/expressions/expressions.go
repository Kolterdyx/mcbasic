package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ExprType string

const (
	BinaryExprType       ExprType = "Binary"
	GroupingExprType     ExprType = "Grouping"
	LiteralExprType      ExprType = "Literal"
	UnaryExprType        ExprType = "Unary"
	VariableExprType     ExprType = "Variable"
	FunctionCallExprType ExprType = "FunctionCall"
	LogicalExprType      ExprType = "Logical"
	SliceExprType        ExprType = "SliceString"
	ListExprType         ExprType = "List"
)

type ExprVisitor interface {
	VisitBinary(b BinaryExpr) interface{}
	VisitGrouping(g GroupingExpr) interface{}
	VisitLiteral(l LiteralExpr) interface{}
	VisitUnary(u UnaryExpr) interface{}
	VisitVariable(v VariableExpr) interface{}
	VisitFunctionCall(f FunctionCallExpr) interface{}
	VisitLogical(l LogicalExpr) interface{}
	VisitSlice(s SliceExpr) interface{}
	VisitList(s ListExpr) interface{}
}

type Expr interface {
	Accept(v ExprVisitor) interface{}
	TType() ExprType
	ReturnType() interfaces.ValueType
}
