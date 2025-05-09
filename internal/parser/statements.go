package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
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
		if p.check(tokens.Equal) || p.check(tokens.BracketOpen) || p.check(tokens.Dot) {
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
			if !p.isListType(varType) {
				return nil, p.error(p.previous(), "Cannot assign empty list to non-list type.")
			}
			if len(list.Elements) == 0 {
				list.ValueType = varType
			}
			initializer = list
		}
		if varType == nil {
			varType = initializer.ReturnType()
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

func (p *Parser) ParseType() (types.ValueType, error) {
	var varType types.ValueType

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
		structStmt, ok := p.structs[structName]
		if !ok {
			return nil, p.error(p.previous(), fmt.Sprintf("Struct '%s' is not defined.", structName))
		}
		varType = structStmt.StructType
	default:
		return nil, p.error(p.peek(), "Expected type.")
	}
	return p.typePostfix(varType)
}

func (p *Parser) typePostfix(varType types.ValueType) (types.ValueType, error) {
	for {
		if !p.match(tokens.BracketOpen) {
			break
		}
		_, err := p.consume(tokens.BracketClose, "Expected ']' after list type.")
		if err != nil {
			return nil, err
		}
		varType = types.NewListType(varType)
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
	if _, err := p.consume(tokens.Equal, "Expected '=' after variable target."); err != nil {
		return nil, err
	}
	valueExpr, err := p.expression()
	if err != nil {
		return nil, err
	}
	targetValueType, err := p.getNestedType(name, accessors)
	if err != nil {
		return nil, err
	}
	if !valueExpr.ReturnType().Equals(targetValueType) {
		return nil, p.error(p.previous(), fmt.Sprintf("Cannot assign %s to %s.", valueExpr.ReturnType().ToString(), targetValueType.ToString()))
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
			var valueType types.ValueType
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
	structType := types.NewStructType(name.Lexeme)
	for !p.check(tokens.BraceClose) && !p.IsAtEnd() {
		fieldName, err := p.consume(tokens.Identifier, "Expected field name.")
		if err != nil {
			return nil, err
		}
		fieldType, err := p.ParseType()
		switch fieldType.(type) {
		case types.PrimitiveTypeStruct:
			switch fieldType {
			case types.IntType:
				structType.SetField(fieldName.Lexeme, types.IntType)
			case types.DoubleType:
				structType.SetField(fieldName.Lexeme, types.DoubleType)
			case types.StringType:
				structType.SetField(fieldName.Lexeme, types.StringType)
			}
		case types.StructTypeStruct:
			structType.SetField(fieldName.Lexeme, fieldType.(types.StructTypeStruct))
		case types.ListTypeStruct:
			structType.SetField(fieldName.Lexeme, fieldType.(types.ListTypeStruct))
		default:
			return nil, p.error(p.previous(), fmt.Sprintf("Invalid field type: %s", fieldType.ToString()))
		}
		if !p.match(tokens.Semicolon) {
			break
		}
	}
	_, err = p.consume(tokens.BraceClose, "Expected '}' after struct fields.")
	if structType.Size() == 0 {
		return nil, p.error(p.peek(), "Struct must have at least one field.")
	}
	stmt := statements.StructDeclarationStmt{
		Name:       name,
		StructType: structType,
	}
	p.structs[name.Lexeme] = stmt
	return stmt, nil
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
