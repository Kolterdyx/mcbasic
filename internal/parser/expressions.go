package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
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

	for p.match(tokens.Slash, tokens.Star) {
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
	if p.match(tokens.Identifier) {
		identifier := p.previous()
		if p.match(tokens.ParenOpen) {
			return p.functionCall(identifier)
		} else {
			return expressions.VariableExpr{Name: identifier, SourceLocation: p.location(), Type: p.getType(identifier)}
		}
	}
	return p.primary()
}

func (p *Parser) functionCall(name tokens.Token) expressions.Expr {
	location := p.location()
	args := make([]expressions.Expr, 0)
	if !p.check(tokens.ParenClose) {
		for {
			args = append(args, p.expression())
			if len(args) >= 255 {
				p.error(p.peek(), "Cannot have more than 255 arguments.")
			}
			if !p.match(tokens.Comma) {
				break
			}
		}
	}
	p.consume(tokens.ParenClose, "Expected ')' after arguments.")
	// Find the function in the current scope
	var returnType expressions.ValueType
	found := false
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
			returnType = f.ReturnType
			found = true
			break
		}
	}
	if !found {
		p.error(name, fmt.Sprintf("Function %s not found.", name.Lexeme))
	}
	return expressions.FunctionCallExpr{Name: name, Arguments: args, SourceLocation: location, Type: returnType}
}

func (p *Parser) primary() expressions.Expr {
	if p.match(tokens.False) {
		return expressions.LiteralExpr{Value: 0, ValueType: expressions.NumberType, SourceLocation: p.location()}
	}
	if p.match(tokens.True) {
		return expressions.LiteralExpr{Value: 1, ValueType: expressions.NumberType, SourceLocation: p.location()}
	}
	if p.match(tokens.Number) {
		return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: expressions.NumberType, SourceLocation: p.location()}
	}
	if p.match(tokens.String) {
		return expressions.LiteralExpr{Value: p.previous().Literal, ValueType: expressions.StringType, SourceLocation: p.location()}
	}
	if p.match(tokens.ParenOpen) {
		expr := p.expression()
		p.consume(tokens.ParenClose, "Expected ')' after expression.")
		return expressions.GroupingExpr{Expression: expr, SourceLocation: p.location()}
	}

	p.error(p.peek(), "Expected expression.")
	return nil
}

func (p *Parser) getType(name tokens.Token) expressions.ValueType {
	// Search the variable in the current scope
	for _, v := range p.variables {
		for _, def := range v {
			if def.Name == name.Lexeme {
				return def.Type
			}
		}
	}
	panic("Variable not found")
}
