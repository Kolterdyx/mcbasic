package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

func (p *Parser) match(tokenTypes ...tokens.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) IsAtEnd() bool {
	return p.current >= len(p.Tokens)
}

func (p *Parser) check(tokenType tokens.TokenType) bool {
	if p.IsAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) peek() tokens.Token {
	if p.IsAtEnd() {
		return p.Tokens[len(p.Tokens)-1]
	}
	return p.Tokens[p.current]
}

func (p *Parser) advance() tokens.Token {
	if !p.IsAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() tokens.Token {
	return p.Tokens[p.current-1]
}

func (p *Parser) error(token tokens.Token, message string) {
	if token.Type == tokens.Eof {
		p.report(token.Line, " at end", message)
	} else {
		p.report(token.Line, " at '"+token.Lexeme+"'", message)
	}
}

func (p *Parser) report(line int, s string, message string) {
	p.HadError = true
	fmt.Printf("[line %d] Error%s: %s\n", line, s, message)
}

func (p *Parser) consume(tokenType tokens.TokenType, message string) tokens.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return tokens.Token{}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.IsAtEnd() {
		if p.previous().Type == tokens.Semicolon {
			return
		}

		switch p.peek().Type {
		case tokens.Let, tokens.Def, tokens.If, tokens.For:
			return
		default:
			p.advance()
		}
	}
}
