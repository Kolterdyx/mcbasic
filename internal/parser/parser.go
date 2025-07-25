package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/types"

	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"slices"
)

type Parser struct {
	tokenSource  []tokens.Token
	errors       []error
	current      int
	definedTypes map[string]types.ValueType
	file         string
}

func NewParser(file string, tokenSource []tokens.Token) *Parser {
	return &Parser{
		tokenSource:  tokenSource,
		errors:       make([]error, 0),
		current:      0,
		definedTypes: make(map[string]types.ValueType),
		file:         file,
	}
}

var allowedTopLevelStatements = []ast.NodeType{
	ast.FunctionDeclarationStatement,
	ast.StructDeclarationStatement,
	ast.VariableDeclarationStatement,
	ast.ImportStatement,
}

func (p *Parser) Parse() (ast.Source, []error) {

	source := make(ast.Source, 0)
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

func (p *Parser) functionDeclarationStatement() (ast.Statement, error) {
	name, err := p.consume(tokens.Identifier, "Expected function name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.ParenOpen, "Expected '(' after function name.")
	if err != nil {
		return nil, err
	}
	parameters := make([]ast.VariableDeclarationStmt, 0)
	parameterTypes := make([]types.ValueType, 0)
	if !p.check(tokens.ParenClose) {
		for {
			if len(parameters) >= 255 {
				return nil, p.error(p.peek(), "Cannot have more than 255 parameters.")
			}
			argName, err := p.consume(tokens.Identifier, "Expected parameter name.")
			if err != nil {
				return nil, err
			}
			argType, err := p.ParseType()
			if err != nil {
				return nil, err
			}
			parameterTypes = append(parameterTypes, argType)
			parameters = append(parameters, ast.VariableDeclarationStmt{
				Name:      argName,
				ValueType: argType,
			})
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	_, err = p.consume(tokens.ParenClose, "Expected ')' after parameters.")
	if err != nil {
		return nil, err
	}
	var returnType types.ValueType = types.VoidType
	if !p.check(tokens.BraceOpen) {
		returnType, err = p.ParseType()
	}
	if returnType == nil {
		return nil, p.error(p.peek(), "Expected return type.")
	}
	if err != nil && !p.check(tokens.BraceOpen) {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return ast.NewFunctionDeclarationStatement(name.Lexeme, parameters, body, returnType), nil
}
