package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type Parser struct {
	Tokens  []tokens.Token
	Headers []interfaces.DatapackHeader

	currentScope string
	current      int

	variables map[string][]interfaces.TypedIdentifier
	functions map[string]interfaces.FuncDef
	Errors    []error
	structs   map[string]statements.StructDeclarationStmt
}

func (p *Parser) Parse() (Program, []error) {
	functions := make(map[string]statements.FunctionDeclarationStmt)
	structs := make(map[string]statements.StructDeclarationStmt)
	p.functions = make(map[string]interfaces.FuncDef)
	for _, funcDef := range GetHeaderFuncDefs(p.Headers) {
		p.functions[funcDef.Name] = funcDef
	}
	p.structs = make(map[string]statements.StructDeclarationStmt)
	p.variables = make(map[string][]interfaces.TypedIdentifier)
	for !p.IsAtEnd() {
		statement, err := p.statement()
		if err != nil {
			p.Errors = append(p.Errors, err)
			continue
		}
		if statement.StmtType() != statements.FunctionDeclarationStmtType && statement.StmtType() != statements.StructDeclarationStmtType {
			p.Errors = append(p.Errors, fmt.Errorf("Only function and struct declarations are allowed at the top level. Found: %s\n", statement.StmtType()))
			continue
		}
		if statement.StmtType() == statements.FunctionDeclarationStmtType {
			funcStmt := statement.(statements.FunctionDeclarationStmt)
			functions[funcStmt.Name.Lexeme] = funcStmt
			p.functions[funcStmt.Name.Lexeme] = interfaces.FuncDef{Name: funcStmt.Name.Lexeme, ReturnType: funcStmt.ReturnType, Args: funcStmt.Parameters}
		}
		if statement.StmtType() == statements.StructDeclarationStmtType {
			structStmt := statement.(statements.StructDeclarationStmt)
			structs[structStmt.Name.Lexeme] = structStmt
		}
	}
	return Program{Functions: functions, Structs: structs}, p.Errors
}
