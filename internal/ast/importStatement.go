package ast

import (
	"fmt"
)

type ImportStmt struct {
	Statement

	Path string
}

func (i ImportStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitImport(i)
}

func (i ImportStmt) Type() NodeType {
	return ImportStatement
}

func (i ImportStmt) ToString() string {
	return fmt.Sprintf("import %s", i.Path)
}
