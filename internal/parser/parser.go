package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/utils"
)

type Parser struct {
	Tokens  []tokens.Token
	Headers []interfaces.DatapackHeader

	currentScope string
	current      int

	variables map[string][]statements.VarDef
	functions []interfaces.FuncDef
	Errors    []error
	structs   map[string]statements.StructDeclarationStmt
}

func (p *Parser) Parse() (Program, []error) {
	var functions []statements.FunctionDeclarationStmt
	var structs []statements.StructDeclarationStmt
	funcDefMap := utils.GetHeaderFuncDefs(p.Headers)
	p.functions = make([]interfaces.FuncDef, 0)
	for _, funcDef := range funcDefMap {
		p.functions = append(p.functions, funcDef)
	}
	p.structs = make(map[string]statements.StructDeclarationStmt)
	p.variables = make(map[string][]statements.VarDef)
	for !p.IsAtEnd() {
		statement, err := p.statement()
		if err != nil {
			p.Errors = append(p.Errors, err)
			continue
		}
		if statement.TType() != statements.FunctionDeclarationStmtType && statement.TType() != statements.StructDeclarationStmtType {
			p.Errors = append(p.Errors, fmt.Errorf("Only function and struct declarations are allowed at the top level. Found: %s\n", statement.TType()))
			continue
		}
		if statement.TType() == statements.FunctionDeclarationStmtType {
			funcStmt := statement.(statements.FunctionDeclarationStmt)
			functions = append(functions, funcStmt)
			p.functions = append(p.functions, interfaces.FuncDef{Name: funcStmt.Name.Lexeme, ReturnType: funcStmt.ReturnType, Args: funcStmt.Parameters})
		}
		if statement.TType() == statements.StructDeclarationStmtType {
			structStmt := statement.(statements.StructDeclarationStmt)
			structs = append(structs, structStmt)
		}
	}
	return Program{Functions: functions, Structs: structs}, p.Errors
}
