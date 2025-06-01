package ast

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type ImportStmt struct {
	interfaces.SourceLocation

	Path        string
	SymbolNames []tokens.Token
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

func (i ImportStmt) GetSourceLocation() interfaces.SourceLocation {
	return i.SourceLocation
}
