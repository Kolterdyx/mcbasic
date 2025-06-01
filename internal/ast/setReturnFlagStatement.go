package ast

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type SetReturnFlagStmt struct {
}

func (s SetReturnFlagStmt) Accept(v StatementVisitor) any {
	return v.VisitSetReturnFlag(s)
}

func (s SetReturnFlagStmt) Type() NodeType {
	return SetReturnFlagStatement
}

func (s SetReturnFlagStmt) ToString() string {
	return "retf;"
}

func (s SetReturnFlagStmt) GetSourceLocation() interfaces.SourceLocation {
	return interfaces.SourceLocation{}
}
