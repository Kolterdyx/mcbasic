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

	for p.match(tokens.BANG_EQUAL, tokens.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() expressions.Expr {
	expr := p.term()

	for p.match(tokens.GREATER, tokens.GREATER_EQUAL, tokens.LESS, tokens.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() expressions.Expr {
	expr := p.factor()

	for p.match(tokens.MINUS, tokens.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() expressions.Expr {
	expr := p.unary()

	for p.match(tokens.SLASH, tokens.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = expressions.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() expressions.Expr {
	if p.match(tokens.BANG, tokens.MINUS) {
		operator := p.previous()
		right := p.unary()
		return expressions.UnaryExpr{Operator: operator, Expression: right}
	}

	return p.primary()
}

func (p *Parser) primary() expressions.Expr {
	if p.match(tokens.FALSE) {
		return expressions.LiteralExpr{Value: false}
	}
	if p.match(tokens.TRUE) {
		return expressions.LiteralExpr{Value: true}
	}
	if p.match(tokens.NUMBER, tokens.STRING) {
		return expressions.LiteralExpr{Value: p.previous().Literal}
	}
	if p.match(tokens.PAREN_OPEN) {
		expr := p.expression()
		p.consume(tokens.PAREN_CLOSE, "Expected ')' after expression.")
		return expressions.GroupingExpr{Expression: expr}
	}
	if p.match(tokens.IDENTIFIER) {
		return expressions.VariableExpr{Name: p.previous()}
	}

	p.error(p.peek(), "Expected expression.")
	return nil
}
