package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type FieldAccessExpr struct {
	Expr
	interfaces.SourceLocation

	Source Expr
	Field  tokens.Token
}

func (v FieldAccessExpr) Accept(visitor ExprVisitor) interfaces.IRCode {
	return visitor.VisitFieldAccess(v)
}

func (v FieldAccessExpr) Type() ast.NodeType {
	return ast.FieldAccessExpression
}

func (v FieldAccessExpr) ToString() string {
	return v.Source.ToString() + "." + v.Field.Lexeme
}
