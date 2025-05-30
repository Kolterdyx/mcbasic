package expressions

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type FunctionCallExpr struct {
	Expr
	interfaces.SourceLocation

	Name      tokens.Token
	Arguments []Expr
	ValueType types.ValueType
}

func (f FunctionCallExpr) Accept(visitor ExprVisitor) interfaces.IRCode {
	return visitor.VisitFunctionCall(f)
}

func (f FunctionCallExpr) Type() ast.NodeType {
	return ast.FunctionCallExpression
}

func (f FunctionCallExpr) ReturnType() types.ValueType {
	return f.ValueType
}

func (f FunctionCallExpr) ToString() string {
	args := ""
	for i, arg := range f.Arguments {
		if i > 0 {
			args += ", "
		}
		args += arg.ToString()
	}
	return f.Name.Lexeme + "(" + args + ")"
}
