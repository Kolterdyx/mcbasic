package parser

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	log "github.com/sirupsen/logrus"
	"strings"
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
	return p.peekCount(0)
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
		p.report(token.Row, token.Col, " at end", message)
	} else {
		p.report(token.Row, token.Col, " at '"+token.Lexeme+"'", message)
	}
}

func (p *Parser) report(line int, pos int, s string, message string) {
	p.HadError = true
	log.Errorf("[Position %d:%d] Error%s: %s\n", line+1, pos+1, s, message)
}

func (p *Parser) consume(tokenType tokens.TokenType, message string) tokens.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	p.error(p.peek(), message)
	return tokens.Token{}
}

func (p *Parser) synchronize() {
	for !p.IsAtEnd() {
		if p.match(tokens.Semicolon, tokens.BraceClose) {
			return
		}
		p.advance()
	}
}

func (p *Parser) stepBack() {
	p.current--
}

func (p *Parser) location() interfaces.SourceLocation {
	return interfaces.SourceLocation{Row: p.previous().Row, Col: p.previous().Col}
}

func (p *Parser) peekCount(offset int) tokens.Token {
	if p.current+offset >= len(p.Tokens) {
		return p.Tokens[len(p.Tokens)-1]
	}
	return p.Tokens[p.current+offset]
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
	for _, f := range p.functions {

		split := strings.Split(f.Name, ":")
		if len(split) > 2 {
			log.Fatalf("Invalid function name: %s", f.Name)
		}
		if len(split) == 2 && split[1] == name.Lexeme || len(split) == 1 && f.Name == name.Lexeme {
			return f.ReturnType
		}
	}
	return ""
}
