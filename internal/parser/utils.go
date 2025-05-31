package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/scanner"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	log "github.com/sirupsen/logrus"
)

func (p *Parser) match(tokenTypes ...interfaces.TokenType) bool {
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
	return p.current >= len(p.tokenSource)
}

func (p *Parser) check(tokenType interfaces.TokenType) bool {
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
	return p.tokenSource[p.current-1]
}

func (p *Parser) error(token tokens.Token, message string) error {
	if token.Type == tokens.Eof {
		return p.report(token.Row, token.Col, " at end", message)
	} else {
		return p.report(token.Row, token.Col, " at '"+token.Lexeme+"'", message)
	}
}

func (p *Parser) report(line int, pos int, s string, message string) error {
	return fmt.Errorf("[Position %d:%d] Syntax error%s: %s\n", line+1, pos+1, s, message)
}

func (p *Parser) consume(tokenType interfaces.TokenType, message string) (tokens.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return tokens.Token{}, p.error(p.peek(), message)
}

func (p *Parser) consumeAny(message string, tokenTypes ...interfaces.TokenType) (tokens.Token, error) {
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
	return interfaces.SourceLocation{File: p.file, Row: p.previous().Row, Col: p.previous().Col}
}

func (p *Parser) peekCount(offset int) tokens.Token {
	if p.current+offset >= len(p.tokenSource) {
		return p.tokenSource[len(p.tokenSource)-1]
	}
	return p.tokenSource[p.current+offset]
}

func (p *Parser) isListType(varType types.ValueType) bool {
	switch varType.(type) {
	case types.ListTypeStruct:
		return true
	default:
		return false
	}
}

func (p *Parser) isStructType(varType types.ValueType) bool {
	switch varType.(type) {
	case types.StructTypeStruct:
		return true
	default:
		return false
	}
}

func parseType(valueType string) (types.ValueType, error) {
	s := scanner.Scanner{}
	tokenSource, errs := s.Scan(valueType)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Errorf("Error parsing type: %s", err)
		}
		return nil, fmt.Errorf("error parsing type: %s", valueType)
	}
	p := Parser{
		tokenSource: tokenSource,
	}
	return p.ParseType()
}

func GetHeaderFuncDefs(header interfaces.DatapackHeader) map[string]interfaces.FunctionDefinition {
	funcDefs := make(map[string]interfaces.FunctionDefinition)
	for _, function := range header.Definitions.Functions {
		funcName := fmt.Sprintf("%s:%s", header.Namespace, function.Name)

		returnType, err := parseType(function.ReturnType)
		if err != nil {
			log.Errorf("Exception parsing function return type: %s", err)
			continue
		}
		_, fn := utils.SplitFunctionName(funcName, "")
		f := interfaces.FunctionDefinition{
			Name:       fn,
			Parameters: make([]interfaces.TypedIdentifier, 0),
			ReturnType: returnType,
		}
		for _, parameter := range function.Args {
			parameterType, err := parseType(parameter.Type)
			if err != nil {
				log.Errorf("Exception parsing function parameter type: %s", err)
				continue
			}
			f.Parameters = append(f.Parameters, interfaces.TypedIdentifier{
				Name: parameter.Name,
				Type: parameterType,
			})
		}
		funcDefs[fn] = f
	}
	return funcDefs
}
