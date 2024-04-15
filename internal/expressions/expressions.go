package expressions

type ExprType string

const (
	BinaryExprType       ExprType = "Binary"
	GroupingExprType     ExprType = "Grouping"
	LiteralExprType      ExprType = "Literal"
	UnaryExprType        ExprType = "Unary"
	VariableExprType     ExprType = "Variable"
	FunctionCallExprType ExprType = "FunctionCall"
	LogicalExprType      ExprType = "Logical"
)

type ExprVisitor interface {
	VisitBinary(b BinaryExpr) interface{}
	VisitGrouping(g GroupingExpr) interface{}
	VisitLiteral(l LiteralExpr) interface{}
	VisitUnary(l UnaryExpr) interface{}
	VisitVariable(l VariableExpr) interface{}
	VisitFunctionCall(l FunctionCallExpr) interface{}
	VisitLogical(l LogicalExpr) interface{}
}

type Expr interface {
	Accept(v ExprVisitor) interface{}
	TType() ExprType
}
