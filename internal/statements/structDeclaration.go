package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type StructDeclarationStmt struct {
	Stmt

	Name     tokens.Token
	Compound *nbt.Compound
}

func (s StructDeclarationStmt) Accept(visitor StmtVisitor) string {
	return visitor.VisitStructDeclaration(s)
}

func (s StructDeclarationStmt) StmtType() StmtType {
	return StructDeclarationStmtType
}
