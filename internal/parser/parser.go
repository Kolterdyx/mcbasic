package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"slices"
)

type Parser struct {
	Tokens  []tokens.Token
	Headers []interfaces.DatapackHeader

	currentScope string
	current      int

	variables map[string][]interfaces.TypedIdentifier
	functions map[string]interfaces.FunctionDefinition
	Errors    []error
	structs   map[string]statements.StructDeclarationStmt
}

var allowedTopLevelStatements = []statements.StmtType{
	statements.FunctionDeclarationStmtType,
	statements.StructDeclarationStmtType,
	statements.ImportStmtType,
}

func (p *Parser) Parse() (Program, []error) {
	functions := make(map[string]statements.FunctionDeclarationStmt)
	structs := make(map[string]statements.StructDeclarationStmt)
	p.functions = make(map[string]interfaces.FunctionDefinition)
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
		if !slices.Contains(allowedTopLevelStatements, statement.StmtType()) {
			p.Errors = append(p.Errors, fmt.Errorf("Found forbidden statement at top level: %s\n", statement.StmtType()))
			continue
		}
		switch statement.StmtType() {
		case statements.FunctionDeclarationStmtType:
			funcStmt := statement.(statements.FunctionDeclarationStmt)
			functions[funcStmt.Name] = funcStmt
			p.functions[funcStmt.Name] = interfaces.FunctionDefinition{Name: funcStmt.Name, ReturnType: funcStmt.ReturnType, Args: funcStmt.Parameters}

		case statements.StructDeclarationStmtType:
			structStmt := statement.(statements.StructDeclarationStmt)
			structs[structStmt.Name.Lexeme] = structStmt
		}
	}
	return Program{Functions: functions, Structs: structs}, p.Errors
}
