package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
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

	var namespaceIdentifier tokens.Token
	var err error
	hasNamespace := false
	if p.match(tokens.Colon) {
		p.stepBack()
		p.stepBack()
	}
	if p.match(tokens.Identifier) {
		identifier := p.previous()
		if p.match(tokens.Colon) {
			namespaceIdentifier = identifier
			hasNamespace = true
			identifier, err = p.consume(tokens.Identifier, "Expected identifier after ':'.")
			if err != nil {
				return nil, err
			}
		}
		identifierType := p.getType(identifier)

		if identifierType == "" {
			return nil, p.error(identifier, "Undeclared identifier")
		}

		switch {
		case p.match(tokens.ParenOpen):
			return p.functionCall(namespaceIdentifier, identifier, hasNamespace)
		case p.match(tokens.BracketOpen):
			return p.slice(expressions.VariableExpr{Name: identifier, SourceLocation: p.location(), Type: identifierType})
		case p.match(tokens.Dot):
			if hasNamespace {
				return nil, p.error(identifier, "Cannot use namespace with field access.")
			}
			return p.fieldAccess(expressions.VariableExpr{Name: identifier, SourceLocation: p.location(), Type: identifierType})
		default:
			return expressions.VariableExpr{Name: identifier, SourceLocation: p.location(), Type: identifierType}, nil
		}
	}
	return p.primary()
}

func (p *Parser) fieldAccess(source expressions.Expr) (expressions.Expr, error) {
	field, err := p.consume(tokens.Identifier, "Expected identifier after '.'")
	if err != nil {
		return nil, err
	}
	fieldType, err := p.getTokenAsValueType(field)
	if err != nil {
		return nil, p.error(field, err.Error())
	}
	fieldAccessExpr := expressions.FieldAccessExpr{
		Source:         source,
		Field:          field,
		SourceLocation: p.location(),
		Type:           fieldType,
	}
	if p.match(tokens.Dot) {
		return p.fieldAccess(fieldAccessExpr)
	}
	return fieldAccessExpr, nil
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
				if arg.ReturnType() != f.Args[i].Type && f.Args[i].Type != expressions.VoidType {
					return nil, p.error(p.peekCount(-(len(f.Args)-i)*2), fmt.Sprintf("Expected %s, got %s.", f.Args[i].Type, arg.ReturnType()))
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
		return expressions.LiteralExpr{Value: 0, ValueType: expressions.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.True) {
		return expressions.LiteralExpr{Value: 1, ValueType: expressions.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.Int) {
		return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: expressions.IntType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.Double) {
		return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: expressions.DoubleType, SourceLocation: p.location()}, nil
	}
	if p.match(tokens.String) {
		if p.match(tokens.BracketOpen) {
			return p.slice(expressions.LiteralExpr{Value: p.peekCount(-2).Literal, ValueType: expressions.StringType, SourceLocation: p.location()})
		} else {
			return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: expressions.StringType, SourceLocation: p.location()}, nil
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
		var contentType = expressions.VoidType
		for !p.check(tokens.BracketClose) {
			if p.match(tokens.Comma) {
				continue
			}
			expr, err := p.expression()
			if err != nil {
				return nil, err
			}
			if contentType == expressions.VoidType {
				contentType = expr.ReturnType()
			} else if contentType != expr.ReturnType() {
				return nil, p.error(p.previous(), fmt.Sprintf("Expected %s, got %s.", contentType, expr.ReturnType()))
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
			ValueType:      p.getListType(contentType),
		}, nil
	}
	errorToken := p.peek()
	p.synchronize()
	return nil, p.error(errorToken, "Expected expression.")
}
