package expressions

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ValueType string

const (
	VoidType   ValueType = ""
	ErrorType  ValueType = "error"
	IntType    ValueType = "int"
	StringType ValueType = "str"
	FixedType  ValueType = "fixed"
)

type LiteralExpr struct {
	Expr
	interfaces.SourceLocation

	Value     interface{}
	ValueType ValueType
}

func (l LiteralExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) TType() ExprType {
	return LiteralExprType
}

func (l LiteralExpr) ReturnType() ValueType {
	return l.ValueType
}
