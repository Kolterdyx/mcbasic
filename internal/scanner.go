package internal

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"unicode"
)

type Scanner struct {
	HadError bool
	source   string
	start    int
	current  int
	tokens   []tokens.Token
	line     int
	pos      int
}

func (s *Scanner) report(line int, column int, message string) {
	log.Errorf("[Position %d:%d] Error: %s\n", line+1, column+1, message)
	s.HadError = true
}

func (s *Scanner) error(line int, message string) {
	s.report(line, s.pos, message)
}

func (s *Scanner) Scan(source string) []tokens.Token {
	s.source = source
	s.tokens = []tokens.Token{}
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()

	switch c {
	case '(':
		s.addToken(tokens.ParenOpen)
	case ')':
		s.addToken(tokens.ParenClose)
	case '{':
		s.addToken(tokens.BraceOpen)
	case '}':
		s.addToken(tokens.BraceClose)
	case ',':
		s.addToken(tokens.Comma)
	case ';':
		s.addToken(tokens.Semicolon)
	case ':':
		s.addToken(tokens.Colon)
	case '-':
		s.addToken(tokens.Minus)
	case '+':
		s.addToken(tokens.Plus)
	case '*':
		s.addToken(tokens.Star)
	case '/':
		s.addToken(tokens.Slash)
	case '%':
		s.addToken(tokens.Percent)
	case '[':
		s.addToken(tokens.BracketOpen)
	case ']':
		s.addToken(tokens.BracketClose)
	case '#':
		s.scanComment()
	case '\n':
		s.line++
		s.pos = 0
		fallthrough
	case ' ', '\r', '\t':
		break
	case '=':
		if s.match('=') {
			s.addToken(tokens.EqualEqual)
			break
		}
		s.addToken(tokens.Equal)
	case '<':
		if s.match('=') {
			s.addToken(tokens.LessEqual)
			break
		}
		s.addToken(tokens.Less)
	case '>':
		if s.match('=') {
			s.addToken(tokens.GreaterEqual)
			break
		}
		s.addToken(tokens.Greater)
	case '!':
		if s.match('=') {
			s.addToken(tokens.BangEqual)
			break
		}
		s.addToken(tokens.Bang)
	case '"':
		s.scanString()
	default:
		if unicode.IsDigit(rune(c)) {
			s.scanNumber()
			break
		}
		if unicode.IsLetter(rune(c)) || c == '_' {
			s.scanIdentifier()
			break
		}
		if c < 32 || c > 126 {
			s.error(s.line, "Unexpected character: "+fmt.Sprintf("'%d'", c))
		} else {
			s.error(s.line, "Unexpected character: "+string(c))
		}
	}

}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.peek() != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peekOffset(offset int) byte {
	if s.current+offset >= len(s.source) {
		return 0
	}
	return s.source[s.current+offset]
}

func (s *Scanner) peek() byte {
	return s.peekOffset(0)
}

func (s *Scanner) advance() byte {
	s.current++
	s.pos++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType tokens.TokenType) {
	s.addTokenWithLiteral(tokenType, "")
}

func (s *Scanner) addTokenWithLiteral(tokenType tokens.TokenType, literal string) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, tokens.Token{Type: tokenType, Lexeme: text, Literal: literal, Line: s.line, Column: s.pos})
}

func (s *Scanner) scanString() {
	stringStartLine := s.line

	for !s.endOfString() && !s.isAtEnd() && s.peek() != '\n' {
		s.advance()
	}
	if s.isAtEnd() || s.peek() == '\n' {
		s.error(s.line, "Unterminated string at line "+fmt.Sprintf("%d", stringStartLine+1))
		return
	}
	s.advance()
	literal := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(tokens.String, s.replaceEscapeSequences(literal))
}

func (s *Scanner) scanComment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
	s.line++
}

func (s *Scanner) scanNumber() {
	for unicode.IsDigit(rune(s.peek())) {
		s.advance()
	}
	if s.peek() == '.' && unicode.IsDigit(rune(s.peekOffset(1))) {
		s.advance()
		for unicode.IsDigit(rune(s.peek())) {
			s.advance()
		}
		num, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
		s.addTokenWithLiteral(tokens.Fixed, strconv.FormatFloat(num, 'f', -1, 64))
	} else if unicode.IsLetter(rune(s.peek())) {
		s.error(s.line, "Unexpected character: "+string(s.peek()))
	} else {
		num, _ := strconv.Atoi(s.source[s.start:s.current])
		s.addTokenWithLiteral(tokens.Int, strconv.Itoa(num))
	}

}

func (s *Scanner) endOfString() bool {
	// Check if the " is escaped by a \. The \ could be escaped by another \, and so on.
	// If the number of \ is odd, the " is escaped.
	// If the number of \ is even, the " is not escaped.
	// If the number of \ is 0, the " is not escaped.

	// Count the number of \ before the current character
	escaped := 0
	for i := s.current - 1; i >= 0; i-- {
		if s.source[i] == '\\' {
			escaped++
		} else {
			break
		}
	}
	return escaped%2 == 0 && s.peek() == '"'
}

func (s *Scanner) replaceEscapeSequences(literal string) string {
	escapeSequences := map[string]string{
		"\\\\": "\\",
		"\\\"": "\"",
	}
	for k, v := range escapeSequences {
		for strings.Contains(literal, k) {
			literal = strings.ReplaceAll(literal, k, v)
		}
	}
	return literal
}

func (s *Scanner) scanIdentifier() {
	for unicode.IsLetter(rune(s.peek())) || unicode.IsDigit(rune(s.peek())) || s.peek() == '_' {
		s.advance()
	}
	text := s.source[s.start:s.current]

	if tokenType, ok := tokens.Keywords[text]; ok {
		s.addToken(tokenType)
	} else {
		s.addTokenWithLiteral(tokens.Identifier, text)
	}
}
