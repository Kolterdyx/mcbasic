package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"slices"
)

type Parser struct {
	tokenSource []tokens.Token
	//Headers []interfaces.DatapackHeader

	symbols       *symbol.Table
	symbolManager *symbol.Manager
	filePath      string

	errors  []error
	current int
}

func NewParser(tokenSource []tokens.Token, manager *symbol.Manager, currentSymbols *symbol.Table) *Parser {
	return &Parser{
		tokenSource:   tokenSource,
		symbols:       currentSymbols,
		symbolManager: manager,
		filePath:      currentSymbols.OriginFile(),
		errors:        make([]error, 0),
		current:       0,
	}
}

var allowedTopLevelStatements = []ast.NodeType{
	ast.FunctionDeclarationStatement,
	ast.StructDeclarationStatement,
	ast.ImportStatement,
}

func (p *Parser) Parse() (*symbol.Table, []statements.Stmt, []error) {

	source := make([]statements.Stmt, 0)
	for !p.IsAtEnd() {
		statement, err := p.statement()
		if err != nil {
			p.errors = append(p.errors, err)
			continue
		}
		if !slices.Contains(allowedTopLevelStatements, statement.Type()) {
			p.errors = append(p.errors, fmt.Errorf("Found forbidden statement at top level: %s\n", statement.Type()))
			continue
		}
		source = append(source, statement)
	}

	return p.symbols, source, p.errors
}
