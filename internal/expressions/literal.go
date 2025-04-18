package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

const (
	// ErrorType  used for error handling
	ErrorType interfaces.ValueType = "error"

	VoidType   interfaces.ValueType = "void"
	IntType    interfaces.ValueType = "int"
	StringType interfaces.ValueType = "str"
	DoubleType interfaces.ValueType = "double"

	ListIntType    interfaces.ValueType = "list<int>"
	ListDoubleType interfaces.ValueType = "list<double>"
	ListStringType interfaces.ValueType = "list<str>"
)

type LiteralExpr struct {
	Expr
	interfaces.SourceLocation

	Value     interface{}
	ValueType interfaces.ValueType
}

func (l LiteralExpr) Accept(v ExprVisitor) string {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) ExprType() ExprType {
	return LiteralExprType
}

func (l LiteralExpr) ReturnType() interfaces.ValueType {
	return l.ValueType
}
