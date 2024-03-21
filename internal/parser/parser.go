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

func (p *Parser) Parse() statements.Stmt {
	return p.block(false)
}
