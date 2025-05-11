package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type SetReturnFlagStmt struct {
	Stmt
}

func (s SetReturnFlagStmt) Accept(v StmtVisitor) interfaces.IRCode {
	return v.VisitSetReturnFlag(s)
}

func (s SetReturnFlagStmt) StmtType() StmtType {
	return SetReturnFlagStmtType
}

func (s SetReturnFlagStmt) ToString() string {
	return "retf;"
}
