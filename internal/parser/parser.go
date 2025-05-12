package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"slices"
)

type Parser struct {
	tokenSource []tokens.Token
	errors      []error
	current     int
}

func NewParser(tokenSource []tokens.Token) *Parser {
	return &Parser{
		tokenSource: tokenSource,
		errors:      make([]error, 0),
		current:     0,
	}
}

var allowedTopLevelStatements = []ast.NodeType{
	ast.FunctionDeclarationStatement,
	ast.StructDeclarationStatement,
	ast.ImportStatement,
}

func (p *Parser) Parse() ([]statements.Stmt, []error) {

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

	return source, p.errors
}
