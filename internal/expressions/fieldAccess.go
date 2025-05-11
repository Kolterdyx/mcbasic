package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type FieldAccessExpr struct {
	Expr
	interfaces.SourceLocation

	Source    Expr
	Field     tokens.Token
	ValueType types.ValueType
}

func (v FieldAccessExpr) Accept(visitor ExprVisitor) interfaces.IRCode {
	return visitor.VisitFieldAccess(v)
}

func (v FieldAccessExpr) Type() ast.NodeType {
	return ast.FieldAccessExpression
}

func (v FieldAccessExpr) ReturnType() types.ValueType {
	return v.ValueType
}

func (v FieldAccessExpr) ToString() string {
	return v.Source.ToString() + "." + v.Field.Lexeme
}
