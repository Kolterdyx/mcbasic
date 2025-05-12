package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type UnaryExpr struct {
	Expr
	interfaces.SourceLocation

	Operator   tokens.Token
	Expression Expr
}

func (u UnaryExpr) Accept(v ExprVisitor) interfaces.IRCode {
	return v.VisitUnary(u)
}

func (u UnaryExpr) Type() ast.NodeType {
	return ast.UnaryExpression
}

func (u UnaryExpr) ReturnType() types.ValueType {
	return types.IntType
}

func (u UnaryExpr) ToString() string {
	return "(" + u.Operator.Lexeme + " " + u.Expression.ToString() + ")"
}
