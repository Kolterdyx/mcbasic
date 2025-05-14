package ast

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

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

	ExpressionStatement          NodeType = "ExpressionStatement"
	VariableDeclarationStatement NodeType = "VariableDeclarationStatement"
	FunctionDeclarationStatement NodeType = "FunctionDeclarationStatement"
	StructDeclarationStatement   NodeType = "StructDeclarationStatement"
	VariableAssignmentStatement  NodeType = "VariableAssignmentStatement"
	BlockStatement               NodeType = "BlockStatement"
	WhileStatement               NodeType = "WhileStatement"
	IfStatement                  NodeType = "IfStatement"
	ReturnStatement              NodeType = "ReturnStatement"
	SetReturnFlagStatement       NodeType = "SetReturnFlagStatement"
	ImportStatement              NodeType = "ImportStatement"

	IndexAccessorType NodeType = "IndexAccessorType"
	FieldAccessorType NodeType = "FieldAccessorType"
)

type Node interface {
	Type() NodeType
	GetSourceLocation() interfaces.SourceLocation
}
