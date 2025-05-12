package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type FunctionCallExpr struct {
	Expr
	interfaces.SourceLocation

	Name      tokens.Token
	Arguments []Expr
}

func (f FunctionCallExpr) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitFunctionCall(f)
}

func (f FunctionCallExpr) Type() NodeType {
	return FunctionCallExpression
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

func (f FunctionCallExpr) GetSourceLocation() interfaces.SourceLocation {
	return f.SourceLocation
}
