package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

func (p *Parser) expression() expressions.Expr {
	return p.or()
}

func (p *Parser) or() expressions.Expr {
	expr := p.and()

	for p.match(tokens.Or) {
		operator := p.previous()
		right := p.and()
		expr = expressions.LogicalExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr
}

func (p *Parser) and() expressions.Expr {
	expr := p.equality()

	for p.match(tokens.And) {
		operator := p.previous()
		right := p.equality()
		expr = expressions.LogicalExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr
}

func (p *Parser) equality() expressions.Expr {
	expr := p.comparison()

	for p.match(tokens.BangEqual, tokens.EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr
}

func (p *Parser) comparison() expressions.Expr {
	expr := p.term()

	for p.match(tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr
}

func (p *Parser) term() expressions.Expr {
	expr := p.factor()

	for p.match(tokens.Minus, tokens.Plus) {
		operator := p.previous()
		right := p.factor()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr
}

func (p *Parser) factor() expressions.Expr {
	expr := p.unary()

	for p.match(tokens.Slash, tokens.Star, tokens.Percent) {
		operator := p.previous()
		right := p.unary()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right, SourceLocation: p.location()}
	}

	return expr
}

func (p *Parser) unary() expressions.Expr {
	if p.match(tokens.Bang, tokens.Minus) {
		operator := p.previous()
		right := p.unary()
		return expressions.UnaryExpr{Operator: operator, Expression: right, SourceLocation: p.location()}
	}

	return p.value()
}

func (p *Parser) value() expressions.Expr {

	var namespaceIdentifier *tokens.Token
	if p.match(tokens.Colon) {
		i := p.previous()
		namespaceIdentifier = &i
	}
	if p.match(tokens.Identifier) {
		identifier := p.previous()
		identifierType := p.getType(identifier)

		if identifierType == "" {
			p.error(identifier, "Undeclared identifier")
			return nil
		}

		if p.match(tokens.ParenOpen) {
			return p.functionCall(namespaceIdentifier, identifier)
		} else if p.match(tokens.BracketOpen) {
			return p.slice(expressions.VariableExpr{Name: identifier, SourceLocation: p.location(), Type: identifierType})
		} else {
			return expressions.VariableExpr{Name: identifier, SourceLocation: p.location(), Type: identifierType}
		}
	}
	return p.primary()
}

func (p *Parser) functionCall(namespace *tokens.Token, name tokens.Token) expressions.Expr {
	location := p.location()
	args := make([]expressions.Expr, 0)
	if !p.check(tokens.ParenClose) {
		for {
			exp := p.expression()
			if exp == nil {
				return nil
			}
			args = append(args, exp)
			if len(args) >= 255 {
				p.error(p.peek(), "Cannot have more than 255 arguments.")
				return nil
			}
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	p.consume(tokens.ParenClose, "Expected ')' after arguments.")
	// Find the function in the current scope
	var funcDef *statements.FuncDef = nil
	for _, f := range p.functions {
		if f.Name == name.Lexeme {
			if len(f.Parameters) != len(args) {
				p.error(name, fmt.Sprintf("Expected %d arguments, got %d.", len(f.Parameters), len(args)))
				return nil
			}
			for i, arg := range args {
				if arg.ReturnType() != f.Parameters[i].Type && f.Parameters[i].Type != expressions.VoidType {
					p.error(p.peekCount(-(len(f.Parameters)-i)*2), fmt.Sprintf("Expected %s, got %s.", f.Parameters[i].Type, arg.ReturnType()))
					return nil
				}
			}
			funcDef = &f
			break
		}
	}
	if funcDef == nil {
		p.error(name, fmt.Sprintf("Function %s not found.", name.Lexeme))
		return nil
	}
	return expressions.FunctionCallExpr{Name: name, Arguments: args, SourceLocation: location, Type: funcDef.ReturnType}
}

func (p *Parser) slice(expr expressions.Expr) expressions.Expr {
	start := p.expression()
	if start == nil {
		return nil
	}
	var end expressions.Expr
	if p.match(tokens.Colon) {
		end = p.expression()
		if end == nil {
			return nil
		}
	}
	p.consume(tokens.BracketClose, "Expected ']' at the end of the slice.")
	return expressions.SliceExpr{StartIndex: start, EndIndex: end, TargetExpr: expr, SourceLocation: p.location()}
}

func (p *Parser) primary() expressions.Expr {
	if p.match(tokens.False) {
		return expressions.LiteralExpr{Value: 0, ValueType: expressions.IntType, SourceLocation: p.location()}
	}
	if p.match(tokens.True) {
		return expressions.LiteralExpr{Value: 1, ValueType: expressions.IntType, SourceLocation: p.location()}
	}
	if p.match(tokens.Int) {
		return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: expressions.IntType, SourceLocation: p.location()}
	}
	if p.match(tokens.Fixed) {
		return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: expressions.FixedType, SourceLocation: p.location()}
	}
	if p.match(tokens.String) {
		if p.match(tokens.BracketOpen) {
			return p.slice(expressions.LiteralExpr{Value: p.peekCount(-2).Literal, ValueType: expressions.StringType, SourceLocation: p.location()})
		} else {
			return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: expressions.StringType, SourceLocation: p.location()}
		}
	}
	if p.match(tokens.ParenOpen) {
		expr := p.expression()
		if expr != nil {
			return nil
		}
		p.consume(tokens.ParenClose, "Expected ')' after expression.")
		return expressions.GroupingExpr{Expression: expr, SourceLocation: p.location()}
	}

	p.error(p.peek(), "Expected expression.")
	p.synchronize()
	return nil
}
