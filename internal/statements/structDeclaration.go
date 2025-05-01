package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type StructDeclarationStmt struct {
	Stmt

	Name       tokens.Token
	StructType types.StructTypeStruct
	Compound   *nbt.Compound
}

func (s StructDeclarationStmt) Accept(visitor StmtVisitor) string {
	return visitor.VisitStructDeclaration(s)
}

func (s StructDeclarationStmt) StmtType() StmtType {
	return StructDeclarationStmtType
}
