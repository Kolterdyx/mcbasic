package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type StructDeclarationStmt struct {
	Statement

	Name       tokens.Token
	StructType types.StructTypeStruct
}

func (s StructDeclarationStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitStructDeclaration(s)
}

func (s StructDeclarationStmt) Type() NodeType {
	return StructDeclarationStatement
}

func (s StructDeclarationStmt) ToString() string {
	return "struct " + s.Name.Lexeme
}
