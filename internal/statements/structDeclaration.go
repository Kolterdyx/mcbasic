package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type StructDeclarationStmt struct {
	Stmt

	Name   tokens.Token
	Fields []interfaces.StructField
}

func (s StructDeclarationStmt) Accept(visitor StmtVisitor) string {
	return visitor.VisitStructDeclaration(s)
}

func (s StructDeclarationStmt) TType() StmtType {
	return StructDeclarationStmtType
}
