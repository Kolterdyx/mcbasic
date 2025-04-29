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

	variables map[string][]statements.VarDef
	functions map[string]interfaces.FuncDef
	Errors    []error
	structs   map[string]statements.StructDeclarationStmt
}

func (p *Parser) Parse() (Program, []error) {
	var functions []statements.FunctionDeclarationStmt
	var structs []statements.StructDeclarationStmt
	p.functions = make(map[string]interfaces.FuncDef, 0)
	for _, funcDef := range GetHeaderFuncDefs(p.Headers) {
		p.functions[funcDef.Name] = funcDef
	}
	p.structs = make(map[string]statements.StructDeclarationStmt)
	p.variables = make(map[string][]statements.VarDef)
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
			functions = append(functions, funcStmt)
			p.functions[funcStmt.Name.Lexeme] = interfaces.FuncDef{Name: funcStmt.Name.Lexeme, ReturnType: funcStmt.ReturnType, Args: funcStmt.Parameters}
		}
		if statement.StmtType() == statements.StructDeclarationStmtType {
			structStmt := statement.(statements.StructDeclarationStmt)
			structs = append(structs, structStmt)
		}
	}
	return Program{Functions: functions, Structs: structs}, p.Errors
}
