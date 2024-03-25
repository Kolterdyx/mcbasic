package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
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
		if statement.Type() != statements.FunctionDeclarationStmtType {
			fmt.Println("Only function declarations are allowed at the top level. Found: ", statement.Type())
		}
		functions = append(functions, statement.(statements.FunctionDeclarationStmt))
	}
	return Program{Functions: functions}
}

func (p *Parser) printStatement() statements.Stmt {
	value := p.expression()
	p.consume(tokens.Semicolon, "Expected ';' after value.")
	return statements.PrintStmt{Expression: value}
}
