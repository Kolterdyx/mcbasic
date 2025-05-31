package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type StructDeclarationStmt struct {
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

func (s StructDeclarationStmt) GetSourceLocation() interfaces.SourceLocation {
	return s.Name.SourceLocation
}
