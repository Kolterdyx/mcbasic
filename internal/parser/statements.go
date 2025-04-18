package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

func (p *Parser) statement() (statements.Stmt, error) {
	switch {
	case p.match(tokens.Let):
		return p.letDeclaration()
	case p.match(tokens.Func):
		return p.functionDeclaration()
	case p.match(tokens.While):
		return p.whileStatement()
	case p.match(tokens.If):
		return p.ifStatement()
	case p.match(tokens.Return):
		return p.returnStatement()
	case p.match(tokens.Struct):
		return p.structDeclaration()
	case p.match(tokens.Identifier):
		if p.check(tokens.Equal) || p.check(tokens.BracketOpen) {
			return p.variableAssignment()
		} else if p.check(tokens.ParenOpen) {
			p.stepBack()
		}
		fallthrough
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) expressionStatement() (statements.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(tokens.Semicolon, "Expected ';' after value."); err != nil {
		return nil, err
	}
	return statements.ExpressionStmt{Expression: value}, nil
}

func (p *Parser) letDeclaration() (statements.Stmt, error) {
	name, err := p.consume(tokens.Identifier, "Expected variable name.")
	if err != nil {
		return nil, err
	}
	var varType interfaces.ValueType
	typeToken, err := p.consumeAny("Expected variable type.", tokens.ValueTypes...)
	if err != nil {
		// Check if the type is a struct
		if p.match(tokens.Identifier) {
			// Check if the struct is defined
			if _, ok := p.structs[p.previous().Lexeme]; !ok {
				return nil, p.error(p.previous(), fmt.Sprintf("Struct '%s' is not defined.", p.previous().Lexeme))
			}
			typeToken = p.previous()
			err = nil
		} else {
			return nil, err
		}
	}
	varType, err = p.getTokenAsValueType(typeToken)
	if err != nil {
		return nil, err
	}
	var initializer expressions.Expr
	if p.match(tokens.Equal) {
		if initializer, err = p.expression(); err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' or '=' after variable declaration.")
	if err != nil {
		return nil, err
	}
	if initializer != nil && initializer.ReturnType() != varType {
		if !(p.isListType(varType) && initializer.ReturnType() == expressions.VoidType) {
			return nil, p.error(p.peekCount(-2), fmt.Sprintf("Cannot assign %s to %s.", initializer.ReturnType(), varType))
		}
	}
	p.variables[p.currentScope] = append(p.variables[p.currentScope], statements.VarDef{
		Name: name.Lexeme,
		Type: varType,
	})
	return statements.VariableDeclarationStmt{
		Name:        name,
		Type:        varType,
		Initializer: initializer,
	}, nil
}

func (p *Parser) variableAssignment() (statements.Stmt, error) {
	name := p.previous()
	var index expressions.Expr
	var err error
	if p.match(tokens.BracketOpen) {
		index, err = p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tokens.BracketClose, "Expected ']' after index.")
		if err != nil {
			return nil, err
		}
		if index.ReturnType() != expressions.IntType {
			return nil, p.error(p.peek(), fmt.Sprintf("Index must be of type %s.", expressions.IntType))
		}
		if !p.isList(name) {
			return nil, p.error(name, fmt.Sprintf("Cannot index type %s.", p.getType(name)))
		}
	}
	_, err = p.consume(tokens.Equal, "Expected '=' after variable name.")
	if err != nil {
		return nil, err
	}
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' after value.")
	if err != nil {
		return nil, err
	}
	return statements.VariableAssignmentStmt{Name: name, Value: value, Index: index}, nil
}

func (p *Parser) functionDeclaration() (statements.Stmt, error) {
	name, err := p.consume(tokens.Identifier, "Expected function name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.ParenOpen, "Expected '(' after function name.")
	if err != nil {
		return nil, err
	}
	parameters := make([]interfaces.FuncArg, 0)
	if !p.check(tokens.ParenClose) {
		for {
			if len(parameters) >= 255 {
				return nil, p.error(p.peek(), "Cannot have more than 255 parameters.")
			}
			argName, err := p.consume(tokens.Identifier, "Expected parameter name.")
			if err != nil {
				return nil, err
			}
			if !p.match(tokens.IntType, tokens.StringType, tokens.DoubleType) {
				err = p.error(p.peek(), "Expected parameter type.")
				p.synchronize()
				return nil, err
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
				return nil, p.error(type_, "Expected parameter type.")
			}
			parameters = append(parameters, interfaces.FuncArg{Name: argName.Lexeme, Type: valueType})
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	_, err = p.consume(tokens.ParenClose, "Expected ')' after parameters.")
	if err != nil {
		return nil, err
	}
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
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return statements.FunctionDeclarationStmt{Name: name, Parameters: parameters, ReturnType: returnType, Body: body}, nil
}

func (p *Parser) structDeclaration() (statements.Stmt, error) {
	name, err := p.consume(tokens.Identifier, "Expected struct name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.BraceOpen, "Expected '{' after struct name.")
	if err != nil {
		return nil, err
	}
	fields := make([]interfaces.StructField, 0)
	for !p.check(tokens.BraceClose) && !p.IsAtEnd() {
		fieldName, err := p.consume(tokens.Identifier, "Expected field name.")
		if err != nil {
			return nil, err
		}
		var fieldType interfaces.ValueType
		if p.match(tokens.ValueTypes...) {
			fieldType, err = p.getTokenAsValueType(p.previous())
			if err != nil {
				return nil, err
			}
		} else {
			// Check if the type is a struct
			if p.match(tokens.Identifier) {
				// Check if the struct is defined
				if _, ok := p.structs[p.previous().Lexeme]; !ok {
					return nil, p.error(p.previous(), fmt.Sprintf("Struct '%s' is not defined.", p.previous().Lexeme))
				}
				fieldType = interfaces.ValueType(p.previous().Lexeme)
			} else {
				return nil, p.error(p.peek(), "Expected field type.")
			}
		}
		fields = append(fields, interfaces.StructField{Name: fieldName.Lexeme, Type: fieldType})
		if !p.match(tokens.Semicolon) {
			break
		}
	}
	_, err = p.consume(tokens.BraceClose, "Expected '}' after struct fields.")
	if len(fields) == 0 {
		return nil, p.error(p.peek(), "Struct must have at least one field.")
	}
	p.structs[name.Lexeme] = statements.StructDeclarationStmt{Name: name, Fields: fields}
	return statements.StructDeclarationStmt{Name: name, Fields: fields}, nil
}

func (p *Parser) block(checkBraces ...bool) (statements.BlockStmt, error) {
	stmts := make([]statements.Stmt, 0)
	if len(checkBraces) == 0 || checkBraces[0] {
		_, err := p.consume(tokens.BraceOpen, "Expected '{' before block.")
		if err != nil {
			return statements.BlockStmt{}, err
		}
	}
	for !p.check(tokens.BraceClose) && !p.IsAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			p.Errors = append(p.Errors, err)
			p.synchronize()
			continue
		}
		stmts = append(stmts, stmt)
	}
	if len(checkBraces) == 0 || checkBraces[0] {
		_, err := p.consume(tokens.BraceClose, "Expected '}' after block.")
		if err != nil {
			return statements.BlockStmt{}, err
		}
	}
	return statements.BlockStmt{Statements: stmts}, nil
}

func (p *Parser) whileStatement() (statements.Stmt, error) {
	_, err := p.consume(tokens.ParenOpen, "Expected '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.ParenClose, "Expected ')' after condition.")
	if err != nil {
		return nil, err
	}
	body, err := p.block()
	if err != nil {
		return nil, err
	}
	return statements.WhileStmt{Condition: condition, Body: body}, nil
}

func (p *Parser) ifStatement() (statements.Stmt, error) {
	_, err := p.consume(tokens.ParenOpen, "Expected '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.ParenClose, "Expected ')' after condition.")
	if err != nil {
		return nil, err
	}
	thenBranch, err := p.block()
	if err != nil {
		return nil, err
	}
	var elseBranch statements.BlockStmt
	if p.match(tokens.Else) {
		if p.match(tokens.If) {
			branch, err := p.ifStatement()
			if err != nil {
				return nil, err
			}
			elseBranch.Statements = append(elseBranch.Statements, branch)
		} else {
			elseBranch, err = p.block()
			if err != nil {
				return nil, err
			}
		}
	}
	return statements.IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser) returnStatement() (statements.Stmt, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' after return statement.")
	if err != nil {
		return nil, err
	}
	return statements.ReturnStmt{Expression: value}, nil
}
