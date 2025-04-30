package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
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
	varType, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	var initializer expressions.Expr
	if p.match(tokens.Equal) {
		if initializer, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if initializer != nil {
		if list, ok := initializer.(expressions.ListExpr); ok {
			// Allow assignment as long as varType is a list
			if _, isList := varType.(types.ListTypeStruct); !isList {
				return nil, p.error(p.previous(), "Cannot assign empty list to non-list type.")
			}
			if len(list.Elements) == 0 {
				list.ValueType = varType
			}
			initializer = list
		}
		if initializer.ReturnType() != varType {
			return nil, p.error(p.previous(), fmt.Sprintf("Cannot assign %s to %s.", initializer.ReturnType().ToString(), varType.ToString()))
		}
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' or '=' after variable declaration.")
	if err != nil {
		return nil, err
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

func (p *Parser) ParseType() (interfaces.ValueType, error) {
	var varType interfaces.ValueType

	// types are as follows:
	// primitive types: int, double, str
	// structs: structName
	// lists: int[], double[], str[], str[][], int[][], double[][], structName[], etc.
	// Lists can be nested
	switch {
	case p.match(tokens.IntType):
		varType = types.IntType
	case p.match(tokens.DoubleType):
		varType = types.DoubleType
	case p.match(tokens.StringType):
		varType = types.StringType
	case p.match(tokens.VoidType):
		varType = types.VoidType
	case p.match(tokens.Identifier):
		structName := p.previous().Lexeme
		if _, ok := p.structs[structName]; !ok {
			return nil, p.error(p.previous(), fmt.Sprintf("Struct '%s' is not defined.", structName))
		}
		varType = types.StructTypeStruct{Name: structName}
	default:
		return nil, p.error(p.peek(), "Expected type.")
	}
	if p.check(tokens.BracketOpen) {
		var listType types.ListTypeStruct
		for p.match(tokens.BracketOpen) {
			if varType == types.VoidType {
				return nil, p.error(p.peek(), "Cannot declare empty list.")
			}
			listType = types.ListTypeStruct{Parent: varType}
			varType = listType
			if !p.match(tokens.BracketClose) {
				return nil, p.error(p.peek(), "Expected ']' after list type.")
			}
		}
		varType = listType
	}
	return varType, nil
}

// variableAssignment parses assignments to variables, list indices, or struct fields
func (p *Parser) variableAssignment() (statements.Stmt, error) {
	name := p.previous()

	// Collect any number of [index] or .field accessors
	var accessors []statements.Accessor
	for {
		if p.match(tokens.BracketOpen) {
			idxExpr, err := p.expression()
			if err != nil {
				return nil, err
			}
			if _, err := p.consume(tokens.BracketClose, "Expected ']' after index"); err != nil {
				return nil, err
			}
			accessors = append(accessors, statements.IndexAccessor{Index: idxExpr})
			continue
		}
		if p.match(tokens.Dot) {
			fieldTok, err := p.consume(tokens.Identifier, "Expected field name after '.'")
			if err != nil {
				return nil, err
			}
			accessors = append(accessors, statements.FieldAccessor{Field: fieldTok})
			continue
		}
		break
	}

	// TODO: Traverse result: var.field[index].etc to make sure the types are compatible

	if _, err := p.consume(tokens.Equal, "Expected '=' after variable target."); err != nil {
		return nil, err
	}
	valueExpr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(tokens.Semicolon, "Expected ';' after assignment."); err != nil {
		return nil, err
	}
	return statements.VariableAssignmentStmt{
		Name:      name,
		Accessors: accessors,
		Value:     valueExpr,
	}, nil
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
				valueType = types.StringType
			case tokens.IntType:
				valueType = types.IntType
			case tokens.DoubleType:
				valueType = types.DoubleType
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
	var returnType interfaces.ValueType = types.VoidType
	if !p.check(tokens.BraceOpen) {
		returnType, err = p.ParseType()
	}
	if returnType == nil {
		return nil, p.error(p.peek(), "Expected return type.")
	}
	if err != nil && !p.check(tokens.BraceOpen) {
		return nil, err
	}
	// Add all parameters to the current scope
	for _, arg := range parameters {
		p.variables[p.currentScope] = append(p.variables[p.currentScope], statements.VarDef{Name: arg.Name, Type: arg.Type})
	}
	p.functions[name.Lexeme] = interfaces.FuncDef{Name: name.Lexeme, Args: parameters, ReturnType: returnType}
	p.currentScope = name.Lexeme
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
	compound := nbt.NewCompound()
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
				fieldType = types.StructTypeStruct{Name: p.previous().Lexeme}
			} else {
				return nil, p.error(p.peek(), "Expected field type.")
			}
		}
		switch fieldType.(type) {
		case types.PrimitiveTypeStruct:
			switch fieldType {
			case types.IntType:
				compound.Set(fieldName.Lexeme, nbt.NewInt(0))
			case types.DoubleType:
				compound.Set(fieldName.Lexeme, nbt.NewDouble(0))
			case types.StringType:
				compound.Set(fieldName.Lexeme, nbt.NewString(""))
			}
		case types.StructTypeStruct:
			structStmt, _ := p.structs[fieldName.Lexeme]
			compound.Set(fieldName.Lexeme, structStmt.Compound)
		case types.ListTypeStruct:
			compound.Set(fieldName.Lexeme, nbt.NewList())
		default:
			return nil, p.error(p.previous(), fmt.Sprintf("Invalid field type: %s", fieldType.ToString()))
		}
		if !p.match(tokens.Semicolon) {
			break
		}
	}
	_, err = p.consume(tokens.BraceClose, "Expected '}' after struct fields.")
	if compound.Size() == 0 {
		return nil, p.error(p.peek(), "Struct must have at least one field.")
	}
	p.structs[name.Lexeme] = statements.StructDeclarationStmt{Name: name, Compound: compound}
	return statements.StructDeclarationStmt{Name: name, Compound: compound}, nil
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
	var expr expressions.Expr = nil
	var err error
	if p.functions[p.currentScope].ReturnType != types.VoidType {
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' after return statement.")
	if err != nil {
		return nil, err
	}
	return statements.ReturnStmt{Expression: expr}, nil
}
