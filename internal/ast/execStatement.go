package ast

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ExecStmt struct {
	SourceLocation interfaces.SourceLocation

	Expression Expr
	ResolvedType
}

func (e ExecStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitExec(e)
}

func (e ExecStmt) Type() NodeType {
	return ExecStatement
}

func (e ExecStmt) ToString() string {
	return "exec " + e.Expression.ToString()
}

func (e ExecStmt) GetSourceLocation() interfaces.SourceLocation {
	return e.SourceLocation
}
