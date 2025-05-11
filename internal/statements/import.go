package statements

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type ImportStmt struct {
	Stmt

	Path string
}

func (i ImportStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitImport(i)
}

func (i ImportStmt) Type() ast.NodeType {
	return ast.ImportStatement
}

func (i ImportStmt) ToString() string {
	return fmt.Sprintf("import %s", i.Path)
}
