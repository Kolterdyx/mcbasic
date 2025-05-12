package frontend

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type CompilationUnit struct {
	FilePath string
	AST      []ast.Statement
	Symbols  *symbol.Table
	Tokens   []tokens.Token
}
