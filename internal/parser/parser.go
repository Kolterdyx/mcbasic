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
	functions []interfaces.FuncDef
}

func (p *Parser) Parse() Program {
	var functions []statements.FunctionDeclarationStmt
	p.functions = make([]interfaces.FuncDef, 0)
	p.variables = make(map[string][]statements.VarDef)
	p.functions = append(p.functions,
		interfaces.FuncDef{Name: "mcb:print", ReturnType: expressions.VoidType, Parameters: []interfaces.FuncArg{{Name: "text", Type: expressions.VoidType}}},
		interfaces.FuncDef{Name: "mcb:log", ReturnType: expressions.VoidType, Parameters: []interfaces.FuncArg{{Name: "text", Type: expressions.VoidType}}},
		interfaces.FuncDef{Name: "mcb:exec", ReturnType: expressions.VoidType, Parameters: []interfaces.FuncArg{{Name: "command", Type: expressions.StringType}}},
		interfaces.FuncDef{Name: "mcb:len", ReturnType: expressions.IntType, Parameters: []interfaces.FuncArg{{Name: "string", Type: expressions.StringType}}},
	)
	for !p.IsAtEnd() {
		statement := p.statement()
		if statement == nil {
			continue
		}
		if statement.TType() != statements.FunctionDeclarationStmtType {
			log.Errorf("Only function declarations are allowed at the top level. Found: %s\n", statement.TType())
			continue
		}
		funcStmt := statement.(statements.FunctionDeclarationStmt)
		functions = append(functions, funcStmt)
		p.functions = append(p.functions, interfaces.FuncDef{Name: funcStmt.Name.Lexeme, ReturnType: funcStmt.ReturnType, Parameters: funcStmt.Parameters})
	}
	return Program{Functions: functions}
}
