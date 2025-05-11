package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type LiteralExpr struct {
	Expr
	interfaces.SourceLocation

	Value     nbt.Value
	ValueType types.ValueType
}

func (l LiteralExpr) Accept(v ExprVisitor) interfaces.IRCode {
	return v.VisitLiteral(l)
}

func (l LiteralExpr) Type() ast.NodeType {
	return ast.LiteralExpression
}

func (l LiteralExpr) ReturnType() types.ValueType {
	return l.ValueType
}

func (l LiteralExpr) ToString() string {
	return l.Value.ToString()
}
