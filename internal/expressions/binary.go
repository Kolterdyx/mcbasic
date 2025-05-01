package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
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

func (b BinaryExpr) ReturnType() types.ValueType {
	switch b.Left.ReturnType() {
	case types.IntType:
		switch b.Right.ReturnType() {
		case types.IntType:
			return types.IntType
		default:
			return types.ErrorType
		}
	case types.StringType:
		return types.StringType
	case types.DoubleType:
		switch b.Right.ReturnType() {
		case types.DoubleType:
			return types.DoubleType
		default:
			return types.ErrorType
		}
	}
	return b.Left.ReturnType()
}
