package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type LogicalExpr struct {
	Left     Expr
	Operator tokens.Token
	Right    Expr

	interfaces.SourceLocation
	Expr
}

func (l LogicalExpr) Accept(v ExprVisitor) interfaces.IRCode {
	return v.VisitLogical(l)
}

func (l LogicalExpr) Type() ast.NodeType {
	return ast.LogicalExpression
}

func (l LogicalExpr) ReturnType() types.ValueType {
	return types.IntType
}

func (l LogicalExpr) ToString() string {
	return "(" + l.Left.ToString() + " " + l.Operator.Lexeme + " " + l.Right.ToString() + ")"
}
