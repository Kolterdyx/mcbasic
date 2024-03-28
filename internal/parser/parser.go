package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	log "github.com/sirupsen/logrus"
)

type Parser struct {
	HadError bool
	current  int
	Tokens   []tokens.Token
}

func (p *Parser) Parse() Program {
	var functions []statements.FunctionDeclarationStmt
	for !p.IsAtEnd() {
		statement := p.statement()
		if statement == nil {
			continue
		}
		if statement.TType() != statements.FunctionDeclarationStmtType {
			log.Errorf("Only function declarations are allowed at the top level. Found: %s\n", statement.TType())
		}
		functions = append(functions, statement.(statements.FunctionDeclarationStmt))
	}
	return Program{Functions: functions}
}

func (p *Parser) stepBack() {
	p.current--
}

func (p *Parser) location() interfaces.SourceLocation {
	return interfaces.SourceLocation{Line: p.previous().Line, Column: p.previous().Column}
}
