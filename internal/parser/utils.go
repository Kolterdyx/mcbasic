package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	log "github.com/sirupsen/logrus"
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
		p.report(token.Line+1, token.Column, " at end", message)
	} else {
		p.report(token.Line+1, token.Column, " at '"+token.Lexeme+"'", message)
	}
	p.synchronize()
}

func (p *Parser) report(line int, pos int, s string, message string) {
	p.HadError = true
	log.Errorf("[Position %d:%d] Error%s: %s\n", line, pos, s, message)
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
		if p.match(tokens.Semicolon, tokens.BraceClose) {
			return
		}
		p.advance()
	}
}
