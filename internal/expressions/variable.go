package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type VariableExpr struct {
	Expr
	interfaces.SourceLocation

	Name      tokens.Token
	ValueType types.ValueType
}

func (v VariableExpr) Accept(visitor ExprVisitor) interfaces.IRCode {
	return visitor.VisitVariable(v)
}

func (v VariableExpr) Type() ast.NodeType {
	return ast.VariableExpression
}

func (v VariableExpr) ReturnType() types.ValueType {
	return v.ValueType
}

func (v VariableExpr) ToString() string {
	return v.Name.Lexeme
}
