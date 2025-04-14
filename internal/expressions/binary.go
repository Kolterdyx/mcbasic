package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type BinaryExpr struct {
	Expr
	interfaces.SourceLocation

	Left     Expr
	Operator tokens.Token
	Right    Expr
}

func (b BinaryExpr) Accept(v ExprVisitor) string {
	return v.VisitBinary(b)
}

func (b BinaryExpr) ExprType() ExprType {
	return BinaryExprType
}

func (b BinaryExpr) ReturnType() interfaces.ValueType {
	switch b.Left.ReturnType() {
	case IntType:
		switch b.Right.ReturnType() {
		case IntType:
			return IntType
		default:
			return ErrorType
		}
	case StringType:
		return StringType
	case DoubleType:
		switch b.Right.ReturnType() {
		case DoubleType:
			return DoubleType
		default:
			return ErrorType
		}
	}
	return b.Left.ReturnType()
}
