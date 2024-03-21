package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

func (p *Parser) expression() expressions.Expr {
	return p.equality()
}

func (p *Parser) equality() expressions.Expr {
	expr := p.comparison()

	for p.match(tokens.BangEqual, tokens.EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() expressions.Expr {
	expr := p.term()

	for p.match(tokens.Greater, tokens.GreaterEqual, tokens.Less, tokens.LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() expressions.Expr {
	expr := p.factor()

	for p.match(tokens.Minus, tokens.Plus) {
		operator := p.previous()
		right := p.factor()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() expressions.Expr {
	expr := p.unary()

	for p.match(tokens.Slash, tokens.Star) {
		operator := p.previous()
		right := p.unary()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() expressions.Expr {
	if p.match(tokens.Bang, tokens.Minus) {
		operator := p.previous()
		right := p.unary()
		return expressions.UnaryExpr{Operator: operator, Expression: right}
	}

	return p.primary()
}

func (p *Parser) primary() expressions.Expr {
	if p.match(tokens.False) {
		return expressions.LiteralExpr{Value: false}
	}
	if p.match(tokens.True) {
		return expressions.LiteralExpr{Value: true}
	}
	if p.match(tokens.Number, tokens.String) {
		return expressions.LiteralExpr{Value: p.previous().Literal}
	}
	if p.match(tokens.ParenOpen) {
		expr := p.expression()
		p.consume(tokens.ParenClose, "Expected ')' after expression.")
		return expressions.GroupingExpr{Expression: expr}
	}
	if p.match(tokens.Identifier) {
		identifier := p.previous()
		if p.match(tokens.ParenOpen) {
			return p.functionCall(identifier)
		} else {
			return expressions.VariableExpr{Name: identifier}
		}
	}

	p.error(p.peek(), "Expected expression.")
	return nil
}
