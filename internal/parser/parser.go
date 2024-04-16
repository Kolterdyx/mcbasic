package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	log "github.com/sirupsen/logrus"
)

type Parser struct {
	HadError bool
	current  int
	Tokens   []tokens.Token

	currentScope string

	variables map[string][]statements.VarDef
	functions []statements.FuncDef
}

func (p *Parser) Parse() Program {
	var functions []statements.FunctionDeclarationStmt
	p.functions = make([]statements.FuncDef, 0)
	p.variables = make(map[string][]statements.VarDef)
	p.functions = append(p.functions,
		statements.FuncDef{Name: "print", ReturnType: expressions.VoidType, Parameters: []statements.FuncArg{{Name: "text", Type: expressions.VoidType}}},
		statements.FuncDef{Name: "exec", ReturnType: expressions.VoidType, Parameters: []statements.FuncArg{{Name: "command", Type: expressions.StringType}}},
	)
	for !p.IsAtEnd() {
		statement := p.statement()
		if statement == nil {
			continue
		}
		if statement.TType() != statements.FunctionDeclarationStmtType {
			log.Errorf("Only function declarations are allowed at the top level. Found: %s\n", statement.TType())
		}
		funcStmt := statement.(statements.FunctionDeclarationStmt)
		functions = append(functions, funcStmt)
		p.functions = append(p.functions, statements.FuncDef{Name: funcStmt.Name.Lexeme, ReturnType: funcStmt.ReturnType, Parameters: funcStmt.Parameters})
	}
	return Program{Functions: functions}
}

func (p *Parser) stepBack() {
	p.current--
}

func (p *Parser) location() interfaces.SourceLocation {
	return interfaces.SourceLocation{Line: p.previous().Line, Column: p.previous().Column}
}

func (p *Parser) peekCount(offset int) tokens.Token {
	if p.current+offset >= len(p.Tokens) {
		return p.Tokens[len(p.Tokens)-1]
	}
	return p.Tokens[p.current+offset]
}
