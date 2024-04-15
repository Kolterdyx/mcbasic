package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

func (p *Parser) statement() statements.Stmt {
	if p.match(tokens.Let) {
		return p.letDeclaration()
	} else if p.match(tokens.Def) {
		return p.functionDeclaration()
	} else if p.match(tokens.While) {
		return p.whileStatement()
	} else if p.match(tokens.If) {
		return p.ifStatement()
	} else if p.match(tokens.Return) {
		return p.returnStatement()
	} else if p.match(tokens.Identifier) {
		if p.check(tokens.Equal) {
			return p.variableAssignment()
		} else if p.check(tokens.ParenOpen) {
			p.stepBack()
			return p.expressionStatement()
		}
	}

	return p.expressionStatement()
}

func (p *Parser) expressionStatement() statements.Stmt {
	value := p.expression()
	p.consume(tokens.Semicolon, "Expected ';' after value.")
	return statements.ExpressionStmt{Expression: value}
}

func (p *Parser) letDeclaration() statements.Stmt {
	name := p.consume(tokens.Identifier, "Expected variable name.")
	var varType tokens.TokenType
	p.consume(tokens.Colon, "Expected type declaration.")
	if p.match(tokens.NumberType, tokens.StringType) {
		varType = p.previous().Type
	} else {
		p.error(p.peek(), "Expected variable type.")
	}
	var initializer expressions.Expr
	if p.match(tokens.Equal) {
		initializer = p.expression()
	}
	p.consume(tokens.Semicolon, "Expected ';' after variable declaration.")
	p.variables[p.currentScope] = append(p.variables[p.currentScope], statements.VarDef{Name: name.Lexeme, Type: varType})
	return statements.VariableDeclarationStmt{Name: name, Type: varType, Initializer: initializer}
}

func (p *Parser) variableAssignment() statements.Stmt {
	name := p.previous()
	p.consume(tokens.Equal, "Expected '=' after variable name.")
	value := p.expression()
	p.consume(tokens.Semicolon, "Expected ';' after value.")
	return statements.VariableAssignmentStmt{Name: name, Value: value}
}

func (p *Parser) functionDeclaration() statements.Stmt {
	name := p.consume(tokens.Identifier, "Expected function name.")
	p.consume(tokens.ParenOpen, "Expected '(' after function name.")
	parameters := make([]statements.FuncArg, 0)
	if !p.check(tokens.ParenClose) {
		for {
			if len(parameters) >= 255 {
				p.error(p.peek(), "Cannot have more than 255 parameters.")
			}
			argName := p.consume(tokens.Identifier, "Expected parameter name.")
			p.consume(tokens.Colon, "Expected type declaration.")
			if !p.match(tokens.NumberType, tokens.StringType) {
				p.error(p.peek(), "Expected parameter type.")
			}
			type_ := p.previous()
			parameters = append(parameters, statements.FuncArg{Name: argName.Lexeme, Type: type_.Type})
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	p.consume(tokens.ParenClose, "Expected ')' after parameters.")
	body := p.block()
	return statements.FunctionDeclarationStmt{Name: name, Parameters: parameters, Body: body}
}

func (p *Parser) block(checkBraces ...bool) statements.BlockStmt {
	stmts := make([]statements.Stmt, 0)
	if len(checkBraces) == 0 || checkBraces[0] {
		p.consume(tokens.BraceOpen, "Expected '{' before block.")
	}
	for !p.check(tokens.BraceClose) && !p.IsAtEnd() {
		stmts = append(stmts, p.statement())
	}
	if len(checkBraces) == 0 || checkBraces[0] {
		p.consume(tokens.BraceClose, "Expected '}' after block.")
	}
	return statements.BlockStmt{Statements: stmts}
}

func (p *Parser) whileStatement() statements.Stmt {
	p.consume(tokens.ParenOpen, "Expected '(' after 'while'.")
	condition := p.expression()
	p.consume(tokens.ParenClose, "Expected ')' after condition.")
	body := p.block()
	return statements.WhileStmt{Condition: condition, Body: body}
}

func (p *Parser) ifStatement() statements.Stmt {
	p.consume(tokens.ParenOpen, "Expected '(' after 'if'.")
	condition := p.expression()
	p.consume(tokens.ParenClose, "Expected ')' after condition.")
	thenBranch := p.block()
	var elseBranch statements.BlockStmt
	if p.match(tokens.Else) {
		if p.match(tokens.If) {
			elseBranch.Statements = append(elseBranch.Statements, p.ifStatement())
		} else {
			elseBranch = p.block()
		}
	}
	return statements.IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) returnStatement() statements.Stmt {
	value := p.expression()
	p.consume(tokens.Semicolon, "Expected ';' after return statement.")
	return statements.ReturnStmt{Expression: value}
}
