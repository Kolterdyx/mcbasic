package frontend

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
)

type CompilationUnit struct {
	FilePath string
	AST      ast.Source
	Symbols  *symbol.Table
}
