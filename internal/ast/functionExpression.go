package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type CallExpr struct {
	interfaces.SourceLocation
	Source    Expr
	Arguments []Expr
	ResolvedType
	resolvedName string
}

func (f *CallExpr) Accept(visitor ExpressionVisitor) any {
	return visitor.VisitCall(f)
}

func (f *CallExpr) Type() NodeType {
	return FunctionCallExpression
}

func (f *CallExpr) ToString() string {
	args := ""
	for i, arg := range f.Arguments {
		if i > 0 {
			args += ", "
		}
		args += arg.ToString()
	}
	return "(" + args + ")"
}

func (f *CallExpr) SetResolvedName(name string) {
	f.resolvedName = name
}

func (f *CallExpr) GetResolvedName() string {
	return f.resolvedName
}

func (f *CallExpr) GetSourceLocation() interfaces.SourceLocation {
	return f.SourceLocation
}
