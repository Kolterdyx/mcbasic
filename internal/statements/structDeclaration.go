package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

// StructField represents a field in a struct declaration.
type StructField struct {
	Name string
	Type interfaces.ValueType
}

type StructDeclarationStmt struct {
	Stmt

	Name   tokens.Token
	Fields []StructField
}

func (s StructDeclarationStmt) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitStructDeclaration(s)
}

func (s StructDeclarationStmt) TType() StmtType {
	return StructDeclarationStmtType
}
