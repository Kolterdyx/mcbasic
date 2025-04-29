package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/scanner"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
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

func (p *Parser) consumeAny(message string, tokenTypes ...tokens.TokenType) (tokens.Token, error) {
	if p.match(tokenTypes...) {
		return p.previous(), nil
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
	return nil
}

func (p *Parser) isListType(varType interfaces.ValueType) bool {
	switch varType.(type) {
	case types.ListTypeStruct:
		return true
	default:
		return false
	}
}

func (p *Parser) getTokenAsValueType(token tokens.Token) (interfaces.ValueType, error) {
	var varType interfaces.ValueType
	var err error
	switch token.Type {
	case tokens.IntType:
		varType = types.IntType
	case tokens.StringType:
		varType = types.StringType
	case tokens.DoubleType:
		varType = types.DoubleType
	case tokens.VoidType:
		varType = types.VoidType
	case tokens.Identifier:
		// structs and lists?
	default:
		return nil, p.error(p.peek(), "Expected variable type.")
	}
	return varType, err
}

func parseType(valueType string) (interfaces.ValueType, error) {
	s := scanner.Scanner{}
	p := Parser{
		Tokens: s.Scan(valueType),
	}
	return p.ParseType()
}

func GetHeaderFuncDefs(headers []interfaces.DatapackHeader) map[string]interfaces.FuncDef {
	funcDefs := make(map[string]interfaces.FuncDef)
	for _, header := range headers {
		log.Debugf("Loading header: %s. Functions: %v", header.Namespace, len(header.Definitions.Functions))
		for _, function := range header.Definitions.Functions {
			funcName := fmt.Sprintf("%s:%s", header.Namespace, function.Name)

			returnType, err := parseType(function.ReturnType)
			if err != nil {
				log.Errorf("Error parsing function return type: %s", err)
				continue
			}
			f := interfaces.FuncDef{
				Name:       funcName,
				Args:       make([]interfaces.FuncArg, 0),
				ReturnType: returnType,
			}
			for _, parameter := range function.Args {
				parameterType, err := parseType(parameter.Type)
				if err != nil {
					log.Errorf("Error parsing function parameter type: %s", err)
					continue
				}
				f.Args = append(f.Args, interfaces.FuncArg{
					Name: parameter.Name,
					Type: parameterType,
				})
			}
			funcDefs[funcName] = f
		}
		log.Debugf("Loaded header: %s. Functions: %v", header.Namespace, len(header.Definitions.Functions))
	}
	return funcDefs
}
