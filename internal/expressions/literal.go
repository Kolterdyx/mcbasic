package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

const (
	VoidType   interfaces.ValueType = "void"
	ErrorType  interfaces.ValueType = "error"
	IntType    interfaces.ValueType = "int"
	StringType interfaces.ValueType = "str"
	DoubleType interfaces.ValueType = "double"
	ListType   interfaces.ValueType = "list"
)

type LiteralExpr struct {
	Expr
	interfaces.SourceLocation

	Value     interface{}
	ValueType interfaces.ValueType
}

func (l LiteralExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) TType() ExprType {
	return LiteralExprType
}

func (l LiteralExpr) ReturnType() interfaces.ValueType {
	return l.ValueType
}
