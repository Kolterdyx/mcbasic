package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

func (p *Parser) statement() (ast.Statement, error) {
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
	case p.match(tokens.Import):
		return p.importStatement()
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

func (p *Parser) importStatement() (ast.Statement, error) {
	importToken := p.previous()
	path, err := p.consume(tokens.String, "Expected path string literal after 'import'.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' after import path.")
	if err != nil {
		return nil, err
	}
	return ast.ImportStmt{Path: path.Lexeme, SourceLocation: importToken.SourceLocation}, nil
}

func (p *Parser) expressionStatement() (ast.Statement, error) {
	value, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(tokens.Semicolon, "Expected ';' after value."); err != nil {
		return nil, err
	}
	return ast.ExpressionStmt{Expression: value}, nil
}

func (p *Parser) letDeclaration() (ast.Statement, error) {
	name, err := p.consume(tokens.Identifier, "Expected variable name.")
	if err != nil {
		return nil, err
	}
	varType, err := p.ParseType()
	if err != nil {
		return nil, err
	}
	var initializer ast.Expr
	if p.match(tokens.Equal) {
		if initializer, err = p.expression(); err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' or '=' after variable declaration.")
	if err != nil {
		return nil, err
	}

	stmt := ast.VariableDeclarationStmt{
		Name:        name,
		Initializer: initializer,
		ValueType:   varType,
	}
	return stmt, nil
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
		typeName := p.previous().Lexeme
		if definedType, ok := p.definedTypes[typeName]; ok {
			varType = definedType
		} else {
			return nil, p.error(p.previous(), fmt.Sprintf("Undefined type: %s", typeName))
		}
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
func (p *Parser) variableAssignment() (ast.Statement, error) {
	name := p.previous()

	// Collect any number of [index] or .field accessors
	var accessors []ast.Accessor
	for {
		if p.match(tokens.BracketOpen) {
			idxExpr, err := p.expression()
			if err != nil {
				return nil, err
			}
			if _, err := p.consume(tokens.BracketClose, "Expected ']' after index"); err != nil {
				return nil, err
			}
			accessors = append(accessors, ast.IndexAccessor{Index: idxExpr})
			continue
		}
		if p.match(tokens.Dot) {
			fieldTok, err := p.consume(tokens.Identifier, "Expected field name after '.'")
			if err != nil {
				return nil, err
			}
			accessors = append(accessors, ast.FieldAccessor{Field: fieldTok})
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
	if _, err := p.consume(tokens.Semicolon, "Expected ';' after assignment."); err != nil {
		return nil, err
	}
	return ast.VariableAssignmentStmt{
		Name:      name,
		Accessors: accessors,
		Value:     valueExpr,
	}, nil
}

func (p *Parser) functionDeclaration() (ast.Statement, error) {
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
	stmt := ast.FunctionDeclarationStmt{Name: name, Parameters: parameters, ReturnType: returnType, Body: body}

	return stmt, nil
}

func (p *Parser) structDeclaration() (ast.Statement, error) {
	name, err := p.consume(tokens.Identifier, "Expected struct name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(tokens.BraceOpen, "Expected '{' after struct name.")
	if err != nil {
		return nil, err
	}
	structType := types.NewStructType(name.Lexeme)
	p.definedTypes[name.Lexeme] = structType
	for !p.check(tokens.BraceClose) && !p.IsAtEnd() {
		fieldName, err := p.consume(tokens.Identifier, "Expected field name.")
		if err != nil {
			return nil, err
		}
		fieldType, err := p.ParseType()
		if err != nil {
			return nil, err
		}
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
	stmt := ast.StructDeclarationStmt{
		Name:       name,
		StructType: structType,
	}
	return stmt, nil
}

func (p *Parser) block(checkBraces ...bool) (ast.BlockStmt, error) {
	stmts := make([]ast.Statement, 0)
	var braceOpen tokens.Token
	var err error
	if len(checkBraces) == 0 || checkBraces[0] {
		braceOpen, err = p.consume(tokens.BraceOpen, "Expected '{' before block.")
		if err != nil {
			return ast.BlockStmt{}, err
		}
	}
	for !p.check(tokens.BraceClose) && !p.IsAtEnd() {
		stmt, err := p.statement()
		if err != nil {
			p.errors = append(p.errors, err)
			p.synchronize()
			continue
		}
		stmts = append(stmts, stmt)
	}
	if len(checkBraces) == 0 || checkBraces[0] {
		_, err := p.consume(tokens.BraceClose, "Expected '}' after block.")
		if err != nil {
			return ast.BlockStmt{}, err
		}
	}
	return ast.BlockStmt{Statements: stmts, SourceLocation: braceOpen.SourceLocation}, nil
}

func (p *Parser) whileStatement() (ast.Statement, error) {
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
	return ast.WhileStmt{Condition: condition, Body: body}, nil
}

func (p *Parser) ifStatement() (ast.Statement, error) {
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
	var elseBranch *ast.BlockStmt = nil
	if p.match(tokens.Else) {
		elseBranch = &ast.BlockStmt{Statements: make([]ast.Statement, 0), SourceLocation: p.previous().SourceLocation}
		if p.match(tokens.If) {
			branch, err := p.ifStatement()
			if err != nil {
				return nil, err
			}
			elseBranch.Statements = append(elseBranch.Statements, branch)
		} else {
			elseB, err := p.block()
			if err != nil {
				return nil, err
			}
			elseBranch = &elseB

		}
	}
	return ast.IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

func (p *Parser) returnStatement() (ast.Statement, error) {
	var expr ast.Expr = nil
	var err error
	ret := p.previous()
	if !p.check(tokens.Semicolon) {
		expr, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokens.Semicolon, "Expected ';' after return statement.")
	if err != nil {
		return nil, err
	}
	return ast.ReturnStmt{Expression: expr, SourceLocation: ret.SourceLocation}, nil
}
