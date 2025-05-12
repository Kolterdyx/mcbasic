package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type VariableExpr struct {
	Expr
	interfaces.SourceLocation

	Name tokens.Token
}

func (v VariableExpr) Accept(visitor ExprVisitor) interfaces.IRCode {
	return visitor.VisitVariable(v)
}

func (v VariableExpr) Type() ast.NodeType {
	return ast.VariableExpression
}

func (v VariableExpr) ToString() string {
	return v.Name.Lexeme
}
