package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

func (p *Parser) statement() statements.Stmt {
	if p.match(tokens.LET) {
		return p.letDeclaration()
	} else if p.match(tokens.DEF) {
		return p.functionDeclaration()
	}
	return p.expressionStatement()
}

func (p *Parser) expressionStatement() statements.Stmt {
	value := p.expression()
	p.consume(tokens.SEMICOLON, "Expected ';' after value.")
	return statements.ExpressionStmt{Expression: value}
}

func (p *Parser) letDeclaration() statements.Stmt {
	name := p.consume(tokens.IDENTIFIER, "Expected variable name.")
	var initializer expressions.Expr
	if p.match(tokens.EQUAL) {
		initializer = p.expression()
	}
	p.consume(tokens.SEMICOLON, "Expected ';' after variable declaration.")
	return statements.VariableDeclarationStmt{Name: name, Initializer: initializer}
}

func (p *Parser) functionDeclaration() statements.Stmt {
	name := p.consume(tokens.IDENTIFIER, "Expected function name.")
	p.consume(tokens.PAREN_OPEN, "Expected '(' after function name.")
	parameters := make([]tokens.Token, 0)
	if !p.check(tokens.PAREN_CLOSE) {
		for {
			if len(parameters) >= 255 {
				p.error(p.peek(), "Cannot have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(tokens.IDENTIFIER, "Expected parameter name."))
			if !p.match(tokens.COMMA) {
				break
			}
		}
	}
	p.consume(tokens.PAREN_CLOSE, "Expected ')' after parameters.")
	p.consume(tokens.BRACE_OPEN, "Expected '{' before function body.")
	body := p.block()
	return statements.FunctionDeclarationStmt{Name: name, Parameters: parameters, Body: body}
}

func (p *Parser) block() []statements.Stmt {
	stmts := make([]statements.Stmt, 0)
	for !p.check(tokens.BRACE_CLOSE) && !p.IsAtEnd() {
		stmts = append(stmts, p.statement())
	}
	p.consume(tokens.BRACE_CLOSE, "Expected '}' after block.")
	return stmts
}
