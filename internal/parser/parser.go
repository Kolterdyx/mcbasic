package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type Parser struct {
	HadError bool
	current  int
	Tokens   []tokens.Token
}

func (p *Parser) Parse() []statements.Stmt {
	stmts := make([]statements.Stmt, 0)
	for !p.IsAtEnd() {
		stmts = append(stmts, p.statement())
	}
	return stmts
}
