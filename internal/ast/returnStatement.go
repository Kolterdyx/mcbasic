package ast

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type ReturnStmt struct {
	interfaces.SourceLocation

	Expression Expr
}

func (s ReturnStmt) Accept(v StatementVisitor) any {
	return v.VisitReturn(s)
}

func (s ReturnStmt) Type() NodeType {
	return ReturnStatement
}

func (s ReturnStmt) ToString() string {
	if s.Expression != nil {
		return "return " + s.Expression.ToString() + ";"
	}
	return "return;"
}

func (s ReturnStmt) GetSourceLocation() interfaces.SourceLocation {
	return s.SourceLocation
}
