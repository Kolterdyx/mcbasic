package ast

type NodeType string

const (
	_                      NodeType = ""
	BinaryExpression       NodeType = "Binary"
	GroupingExpression     NodeType = "Grouping"
	LiteralExpression      NodeType = "Literal"
	UnaryExpression        NodeType = "Unary"
	FieldAccessExpression  NodeType = "FieldAccess"
	VariableExpression     NodeType = "Variable"
	FunctionCallExpression NodeType = "FunctionCall"
	LogicalExpression      NodeType = "Logical"
	SliceExpression        NodeType = "SliceString"
	ListExpression         NodeType = "List"
	StructExpression       NodeType = "Struct"
)

type Node interface {
	Type() NodeType
}
