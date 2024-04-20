package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
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
		statements.FuncDef{Name: "len", ReturnType: expressions.IntType, Parameters: []statements.FuncArg{{Name: "string", Type: expressions.StringType}}},
	)
	for !p.IsAtEnd() {
		statement := p.statement()
		if statement == nil {
			fmt.Println("Error in statement")
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
