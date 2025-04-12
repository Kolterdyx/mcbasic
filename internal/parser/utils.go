package parser

import (
	"fmt"
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

func (p *Parser) error(token tokens.Token, message string) error {
	if token.Type == tokens.Eof {
		return p.report(token.Row, token.Col, " at end", message)
	} else {
		return p.report(token.Row, token.Col, " at '"+token.Lexeme+"'", message)
	}
}

func (p *Parser) report(line int, pos int, s string, message string) error {
	return fmt.Errorf("[Position %d:%d] Error%s: %s\n", line+1, pos+1, s, message)
}

func (p *Parser) consume(tokenType tokens.TokenType, message string) (tokens.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return tokens.Token{}, p.error(p.peek(), message)
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

func (p *Parser) getType(name tokens.Token) interfaces.ValueType {
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

func (p *Parser) isList(name tokens.Token) bool {
	for _, v := range p.variables {
		for _, def := range v {
			if def.Name == name.Lexeme {
				return def.Type == expressions.ListIntType ||
					def.Type == expressions.ListDoubleType ||
					def.Type == expressions.ListStringType
			}
		}
	}
	for _, f := range p.functions {

		split := strings.Split(f.Name, ":")
		if len(split) > 2 {
			log.Fatalf("Invalid function name: %s", f.Name)
		}
		if len(split) == 2 && split[1] == name.Lexeme || len(split) == 1 && f.Name == name.Lexeme {
			return f.ReturnType == expressions.ListIntType ||
				f.ReturnType == expressions.ListDoubleType ||
				f.ReturnType == expressions.ListStringType
		}
	}
	return false
}

func (p *Parser) isListType(varType interfaces.ValueType) bool {
	switch varType {
	case expressions.ListIntType:
		return true
	case expressions.ListStringType:
		return true
	case expressions.ListDoubleType:
		return true
	default:
		return false
	}
}

func (p *Parser) getListType(valueType interfaces.ValueType) interfaces.ValueType {
	switch valueType {
	case expressions.IntType:
		return expressions.ListIntType
	case expressions.StringType:
		return expressions.ListStringType
	case expressions.DoubleType:
		return expressions.ListDoubleType
	case expressions.VoidType:
		return expressions.VoidType
	default:
		log.Fatalf("Unsupported type for list: %s", valueType)
	}
	return ""
}

func (p *Parser) getListValueType(valueType interfaces.ValueType) interfaces.ValueType {
	switch valueType {
	case expressions.ListIntType:
		return expressions.IntType
	case expressions.ListStringType:
		return expressions.StringType
	case expressions.ListDoubleType:
		return expressions.DoubleType
	default:
		log.Fatalf("Invalid list type: %s", valueType)
	}
	return ""
}
