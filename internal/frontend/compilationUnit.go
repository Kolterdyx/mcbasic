package frontend

import (
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type CompilationUnit struct {
	FilePath string
	AST      []statements.Stmt
	Symbols  *symbol.Table
	Tokens   []tokens.Token
}
