package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

func (p *Parser) statement() statements.Stmt {
	if p.match(tokens.Let) {
		return p.letDeclaration()
	} else if p.match(tokens.Func) {
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
	if value == nil {
		return nil
	}
	p.consume(tokens.Semicolon, "Expected ';' after value.")
	return statements.ExpressionStmt{Expression: value}
}

func (p *Parser) letDeclaration() statements.Stmt {
	name := p.consume(tokens.Identifier, "Expected variable name.")
	var varType interfaces.ValueType
	if p.match(tokens.IntType) {
		varType = expressions.IntType
	} else if p.match(tokens.StringType) {
		varType = expressions.StringType
	} else if p.match(tokens.DoubleType) {
		varType = expressions.DoubleType
	} else {
		p.error(p.peek(), "Expected variable type.")
		return nil
	}
	var initializer expressions.Expr
	if p.match(tokens.Equal) {
		initializer = p.expression()
		if initializer == nil {
			return nil
		}
	}
	p.consume(tokens.Semicolon, "Expected ';' after variable declaration.")
	if initializer != nil && initializer.ReturnType() != varType {
		p.error(p.peekCount(-2), fmt.Sprintf("Cannot assign %s to %s.", initializer.ReturnType(), varType))
		return nil
	}
	p.variables[p.currentScope] = append(p.variables[p.currentScope], statements.VarDef{Name: name.Lexeme, Type: varType})
	return statements.VariableDeclarationStmt{Name: name, Type: varType, Initializer: initializer}
}

func (p *Parser) variableAssignment() statements.Stmt {
	name := p.previous()
	p.consume(tokens.Equal, "Expected '=' after variable name.")
	value := p.expression()
	if value == nil {
		return nil
	}
	p.consume(tokens.Semicolon, "Expected ';' after value.")
	return statements.VariableAssignmentStmt{Name: name, Value: value}
}

func (p *Parser) functionDeclaration() statements.Stmt {
	name := p.consume(tokens.Identifier, "Expected function name.")
	p.consume(tokens.ParenOpen, "Expected '(' after function name.")
	parameters := make([]interfaces.FuncArg, 0)
	if !p.check(tokens.ParenClose) {
		for {
			if len(parameters) >= 255 {
				p.error(p.peek(), "Cannot have more than 255 parameters.")
				return nil
			}
			argName := p.consume(tokens.Identifier, "Expected parameter name.")
			if !p.match(tokens.IntType, tokens.StringType, tokens.DoubleType) {
				p.error(p.peek(), "Expected parameter type.")
				p.synchronize()
				return nil
			}
			type_ := p.previous()
			var valueType interfaces.ValueType
			switch type_.Type {
			case tokens.StringType:
				valueType = expressions.StringType
			case tokens.IntType:
				valueType = expressions.IntType
			case tokens.DoubleType:
				valueType = expressions.DoubleType
			default:
				p.error(type_, "Expected parameter type.")
				return nil
			}
			parameters = append(parameters, interfaces.FuncArg{Name: argName.Lexeme, Type: valueType})
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	p.consume(tokens.ParenClose, "Expected ')' after parameters.")
	returnType := expressions.VoidType
	if p.match(tokens.IntType) {
		returnType = expressions.IntType
	} else if p.match(tokens.StringType) {
		returnType = expressions.StringType
	} else if p.match(tokens.DoubleType) {
		returnType = expressions.DoubleType
	}

	// Add all parameters to the current scope
	for _, arg := range parameters {
		p.variables[p.currentScope] = append(p.variables[p.currentScope], statements.VarDef{Name: arg.Name, Type: arg.Type})
	}
	p.functions = append(p.functions, interfaces.FuncDef{Name: name.Lexeme, Args: parameters, ReturnType: returnType})
	body := p.block()
	return statements.FunctionDeclarationStmt{Name: name, Parameters: parameters, ReturnType: returnType, Body: body}
}

func (p *Parser) block(checkBraces ...bool) statements.BlockStmt {
	stmts := make([]statements.Stmt, 0)
	if len(checkBraces) == 0 || checkBraces[0] {
		p.consume(tokens.BraceOpen, "Expected '{' before block.")
	}
	for !p.check(tokens.BraceClose) && !p.IsAtEnd() {
		stmt := p.statement()
		if stmt == nil {
			p.synchronize()
			continue
		}
		stmts = append(stmts, stmt)
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
