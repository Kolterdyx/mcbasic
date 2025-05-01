package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

func (p *Parser) expression() (expressions.Expr, error) {
	return p.or()
}

func (p *Parser) or() (expressions.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Or) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = expressions.LogicalExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) and() (expressions.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.And) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = expressions.LogicalExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) equality() (expressions.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.BangEqual, tokens.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) comparison() (expressions.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) term() (expressions.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Minus, tokens.Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) factor() (expressions.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(tokens.Slash, tokens.Star, tokens.Percent) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr, nil
}

func (p *Parser) unary() (expressions.Expr, error) {
	if p.match(tokens.Bang, tokens.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return expressions.UnaryExpr{Operator: operator, Expression: right, SourceLocation: p.location()}, nil
	}

	return p.value()
}

func (p *Parser) value() (expressions.Expr, error) {
	// 1) Parse the “primary” thing (identifier, literal, function call, grouping, list literal…)
	expr, err := p.baseValue()
	if err != nil {
		return nil, err
	}
	// 2) Repeatedly eat [] or .field, nesting the expression:
	return p.postfix(expr)
}

// baseValue parses an atomic expression with optional namespace and function call.
func (p *Parser) baseValue() (expressions.Expr, error) {
	var namespaceToken tokens.Token
	var hasNamespace bool
	var err error
	if p.match(tokens.Colon) {
		p.stepBack()
		p.stepBack()
	}

	// First, try to parse an identifier (possibly namespaced)
	if p.match(tokens.Identifier) {
		// Save the identifier (might be namespace)
		identifier := p.previous()

		// If a colon follows, treat previous as namespace
		if p.match(tokens.Colon) {
			namespaceToken = identifier
			hasNamespace = true
			// Next token must be the actual identifier
			identifier, err = p.consume(tokens.Identifier, "Expected identifier after ':'")
			if err != nil {
				return nil, err
			}
		}

		// Lookup the declared type
		identifierType := p.getType(identifier)
		if identifierType == nil {
			return nil, p.error(identifier, "Undeclared identifier")
		}

		// If this is a function call, delegate to functionCall (handles namespace)
		if p.match(tokens.ParenOpen) {
			if _, ok := p.structs[identifier.Lexeme]; ok {
				//return p.structLiteral(identifier, structStmt.StructType)
			} else {
				return p.functionCall(namespaceToken, identifier, hasNamespace)
			}
		}

		// Otherwise, it's a variable reference.  If namespaced, prefix the lexeme.
		nameToken := identifier
		if hasNamespace {
			nameToken = tokens.Token{
				Type:           identifier.Type,
				Lexeme:         fmt.Sprintf("%s:%s", namespaceToken.Lexeme, identifier.Lexeme),
				Literal:        identifier.Literal,
				SourceLocation: identifier.SourceLocation,
			}
		}

		return expressions.VariableExpr{
			Name:           nameToken,
			SourceLocation: p.location(),
			Type:           identifierType,
		}, nil
	}

	// Fallback to literals, grouping, list literals, etc.
	return p.primary()
}

// postfix wraps an Expr in as many [index] or .field as you find:
func (p *Parser) postfix(expr expressions.Expr) (expressions.Expr, error) {
	// If it’s “[”, parse a slice/index, then recurse:
	switch returnType := expr.ReturnType().(type) {
	case types.ListTypeStruct, types.PrimitiveTypeStruct:
		if returnType != types.StringType {
			break
		}
		if p.match(tokens.BracketOpen) {
			// read either slice or single index
			start, err := p.expression()
			if err != nil {
				return nil, err
			}
			var end expressions.Expr
			if p.match(tokens.Colon) {
				end, err = p.expression()
				if err != nil {
					return nil, err
				}
			}
			if _, err = p.consume(tokens.BracketClose, "Expected ']'"); err != nil {
				return nil, err
			}
			expr = expressions.SliceExpr{
				TargetExpr:     expr,
				StartIndex:     start,
				EndIndex:       end,
				SourceLocation: p.location(),
			}
			// keep consuming more postfix:
			return p.postfix(expr)
		}
	// If it’s “.”, parse a field access, then recurse:
	case types.StructTypeStruct:
		if p.match(tokens.Dot) {
			fieldTok, err := p.consume(tokens.Identifier, "Expected field name after '.'")
			if err != nil {
				return nil, err
			}
			fieldTokenType, err := p.getTokenAsValueType(fieldTok)
			if err != nil {
				return nil, err
			}
			expr = expressions.FieldAccessExpr{
				Source:         expr,
				Field:          fieldTok,
				SourceLocation: p.location(),
				Type:           fieldTokenType,
			}
			return p.postfix(expr)
		}
	}
	return expr, nil
}

func (p *Parser) functionCall(namespace tokens.Token, name tokens.Token, hasNamespace bool) (expressions.Expr, error) {
	location := p.location()
	args := make([]expressions.Expr, 0)
	if !p.check(tokens.ParenClose) {
		for {
			exp, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, exp)
			if len(args) >= 255 {
				return nil, p.error(p.peek(), "Cannot have more than 255 arguments.")
			}
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	_, err := p.consume(tokens.ParenClose, "Expected ')' after arguments.")
	if err != nil {
		return nil, err
	}
	// Find the function in the current scope

	lexeme := name.Lexeme
	if hasNamespace {
		lexeme = fmt.Sprintf("%s:%s", namespace.Lexeme, name.Lexeme)
	}

	var funcDef *interfaces.FuncDef = nil
	for _, f := range p.functions {
		if f.Name == lexeme {
			if len(f.Args) != len(args) {
				return nil, p.error(name, fmt.Sprintf("Expected %d arguments, got %d.", len(f.Args), len(args)))
			}
			for i, arg := range args {
				if arg.ReturnType() != f.Args[i].Type && f.Args[i].Type != types.VoidType {
					return nil, p.error(p.peekCount(-(len(f.Args)-i)*2), fmt.Sprintf("Expected %s, got %s.", f.Args[i].Type.ToString(), arg.ReturnType().ToString()))
				}
			}
			funcDef = &f
			break
		}
	}
	if funcDef == nil {
		return nil, p.error(name, fmt.Sprintf("Function %s not found.", name.Lexeme))
	}

	return expressions.FunctionCallExpr{Name: tokens.Token{
		Type:           tokens.Identifier,
		Lexeme:         lexeme,
		Literal:        name.Literal,
		SourceLocation: name.SourceLocation,
	}, Arguments: args, SourceLocation: location, Type: funcDef.ReturnType}, nil
}

func (p *Parser) structLiteral(name tokens.Token, structType types.StructTypeStruct) (expressions.Expr, error) {
	var args []expressions.Expr
	if !p.check(tokens.ParenClose) {
		for {
			exp, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, exp)
			if len(args) >= 255 {
				return nil, p.error(p.peek(), "Cannot have more than 255 arguments.")
			}
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	_, err := p.consume(tokens.ParenClose, "Expected ')' after arguments.")
	if err != nil {
		return nil, err
	}

	if len(args) != structType.Size() {
		return nil, p.error(name, fmt.Sprintf("Expected %d arguments, got %d.", structType.Size(), len(args)))
	}

	fieldNames := structType.GetFieldNames()
	for i, arg := range args {
		fieldType, _ := structType.GetField(fieldNames[i])
		if !arg.ReturnType().Equals(fieldType) && !fieldType.Equals(types.VoidType) {
			return nil, p.error(p.peekCount(-(structType.Size()-i)*2), fmt.Sprintf("Expected %s, got %s.", fieldType.ToString(), arg.ReturnType().ToString()))
		}
	}

	return expressions.StructExpr{Args: args, StructType: structType}, nil
}

func (p *Parser) slice(expr expressions.Expr) (expressions.Expr, error) {
	start, err := p.expression()
	if err != nil {
		return nil, err
	}
	var end expressions.Expr
	if p.match(tokens.Colon) {
		end, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(tokens.BracketClose, "Expected ']' at the end of the slice.")
	if err != nil {
		return nil, err
	}
	return expressions.SliceExpr{StartIndex: start, EndIndex: end, TargetExpr: expr, SourceLocation: p.location()}, nil
}

func (p *Parser) primary() (expressions.Expr, error) {
	if p.match(tokens.False) {
		return expressions.LiteralExpr{Value: "0", ValueType: types.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.True) {
		return expressions.LiteralExpr{Value: "1", ValueType: types.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.Int) {
		return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: types.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.Double) {
		return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: types.DoubleType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.String) {
		if p.match(tokens.BracketOpen) {
			return p.slice(expressions.LiteralExpr{Value: p.peekCount(-2).Literal, ValueType: types.StringType, SourceLocation: p.location()})
		} else {
			return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: types.StringType, SourceLocation: p.location()}, nil
		}
	}
	if p.match(tokens.ParenOpen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(tokens.ParenClose, "Expected ')' after expression.")
		if err != nil {
			return nil, err
		}
		return expressions.GroupingExpr{Expression: expr, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.BracketOpen) {
		var elems []expressions.Expr
		var contentType types.ValueType = types.VoidType
		for !p.check(tokens.BracketClose) {
			if p.match(tokens.Comma) {
				continue
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			if contentType == types.VoidType {
				contentType = expr.ReturnType()
			} else if contentType != expr.ReturnType() {
				return nil, p.error(p.previous(), fmt.Sprintf("Expected %s, got %s.", contentType.ToString(), expr.ReturnType().ToString()))
			}
			elems = append(elems, expr)
			if p.check(tokens.BracketClose) {
				break
			}
			_, err = p.consume(tokens.Comma, "Expected ',' after expression.")
			if err != nil {
				return nil, err
			}
		}
		_, err := p.consume(tokens.BracketClose, "Expected ']' after expression.")
		if err != nil {
			return nil, err
		}
		return expressions.ListExpr{
			Elements:       elems,
			SourceLocation: p.location(),
			ValueType:      types.NewListType(contentType),
		}, nil
	}
	errorToken := p.peek()
	p.synchronize()
	return nil, p.error(errorToken, "Expected expression.")
}
