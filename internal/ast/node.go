package ast

type NodeType string

const (
	_                      NodeType = ""
	BinaryExpression       NodeType = "BinaryExpression"
	GroupingExpression     NodeType = "GroupingExpression"
	LiteralExpression      NodeType = "LiteralExpression"
	UnaryExpression        NodeType = "UnaryExpression"
	FieldAccessExpression  NodeType = "FieldAccessExpression"
	VariableExpression     NodeType = "VariableExpression"
	FunctionCallExpression NodeType = "FunctionCallExpression"
	LogicalExpression      NodeType = "LogicalExpression"
	SliceExpression        NodeType = "SliceExpression"
	ListExpression         NodeType = "ListExpression"
	StructExpression       NodeType = "StructExpression"

	ExpressionStatement          NodeType = "ExpressionStatement"
	VariableDeclarationStatement NodeType = "VariableDeclarationStatement"
	FunctionDeclarationStatement NodeType = "FunctionDeclarationStatement"
	StructDeclarationStatement   NodeType = "StructDeclarationStatement"
	VariableAssignmentStatement  NodeType = "VariableAssignmentStatement"
	BlockStatement               NodeType = "BlockStatement"
	WhileStatement               NodeType = "WhileStatement"
	IfStatement                  NodeType = "IfStatement"
	ReturnStatement              NodeType = "ReturnStatement"
	ScoreStatement               NodeType = "ScoreStatement"
	SetReturnFlagStatement       NodeType = "SetReturnFlagStatement"
	ImportStatement              NodeType = "ImportStatement"
)

type Node interface {
	Type() NodeType
}
