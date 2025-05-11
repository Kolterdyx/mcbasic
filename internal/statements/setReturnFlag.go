package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type SetReturnFlagStmt struct {
	Stmt
}

func (s SetReturnFlagStmt) Accept(v StmtVisitor) interfaces.IRCode {
	return v.VisitSetReturnFlag(s)
}

func (s SetReturnFlagStmt) Type() ast.NodeType {
	return ast.SetReturnFlagStatement
}

func (s SetReturnFlagStmt) ToString() string {
	return "retf;"
}
