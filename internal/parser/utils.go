package parser

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/scanner"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
	"reflect"
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

func (p *Parser) getType(name tokens.Token) types.ValueType {
	// Search the variable in the current scope
	for _, varDef := range p.variables[p.currentScope] {
		if varDef.Name == name.Lexeme {
			return varDef.Type
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
	if structStmt, ok := p.structs[name.Lexeme]; ok {
		return structStmt.StructType
	}
	return nil
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
		log.Debugf("isStructType: %+v", reflect.TypeOf(varType))
		return false
	}
}

func (p *Parser) getTokenAsValueType(token tokens.Token) (types.ValueType, error) {
	var varType types.ValueType = types.ErrorType
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

// getNestedType traverses the accessors to find the type at the end
func (p *Parser) getNestedType(name tokens.Token, accessors []statements.Accessor) (types.ValueType, error) {
	varType := p.getType(name)
	if varType == nil {
		return nil, p.error(name, "Undeclared identifier")
	}
	var err error
	accessPath := name.Lexeme
	for _, accessor := range accessors {
		accessPath += accessor.ToString()
		switch accessor.(type) {
		case statements.IndexAccessor:
			if p.isListType(varType) {
				varType = varType.(types.ListTypeStruct).ContentType
			} else {
				return nil, p.error(p.peek(), "Expected list type.")
			}
		case statements.FieldAccessor:
			fieldAccessor := accessor.(statements.FieldAccessor)
			if p.isStructType(varType) {
				vtype, ok := varType.(types.StructTypeStruct).GetField(fieldAccessor.Field.Lexeme)
				if !ok {
					return nil, p.error(fieldAccessor.Field, fmt.Sprintf("Unknown field: %s", fieldAccessor.Field.Lexeme))
				}
				varType = vtype
			} else {
				return nil, p.error(p.peek(), "Expected struct type.")
			}
		default:
			return nil, p.error(p.peek(), "Unknown accessor type.")
		}
	}
	if varType == nil {
		return nil, p.error(name, fmt.Sprintf("Unknown variable type: %s", name.Lexeme))
	}
	return varType, err
}

func parseType(valueType string) (types.ValueType, error) {
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
				Args:       make([]interfaces.TypedIdentifier, 0),
				ReturnType: returnType,
			}
			for _, parameter := range function.Args {
				parameterType, err := parseType(parameter.Type)
				if err != nil {
					log.Errorf("Error parsing function parameter type: %s", err)
					continue
				}
				f.Args = append(f.Args, interfaces.TypedIdentifier{
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
