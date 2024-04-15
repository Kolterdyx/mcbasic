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

func (b BinaryExpr) Accept(v ExprVisitor) interface{} {
	return v.VisitBinary(b)
}

func (b BinaryExpr) TType() ExprType {
	return BinaryExprType
}
