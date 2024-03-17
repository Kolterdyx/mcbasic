package internal

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Scanner struct {
	HadError bool
	source   string
	start    int
	current  int
	tokens   []Token
	line     int
}

func (s *Scanner) report(line int, where string, message string) {
	fmt.Printf("[line %d] Error %s: %s\n", line, where, message)
	s.HadError = true
}

func (s *Scanner) error(line int, message string) {
	s.report(line, "", message)
}

func (s *Scanner) Scan(source string) []Token {
	s.source = source
	s.tokens = []Token{}
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
		s.addToken(PAREN_OPEN)
	case ')':
		s.addToken(PAREN_CLOSE)
	case '{':
		s.addToken(BRACE_OPEN)
	case '}':
		s.addToken(BRACE_CLOSE)
	case ',':
		s.addToken(COMMA)
	case '-':
		s.addToken(MINUS)
	case '+':
		s.addToken(PLUS)
	case '*':
		s.addToken(STAR)
	case '/':
		s.addToken(SLASH)
	case '%':
		s.addToken(PERCENT)
	case '#':
		s.scanComment()
	case '\n':
		s.line++
		fallthrough
	case ' ':
		fallthrough
	case '\r':
		fallthrough
	case '\t':
		break
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		}
		s.addToken(EQUAL)
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		}
		s.addToken(LESS)
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		}
		s.addToken(GREATER)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL)
		}
		s.addToken(BANG)
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
			s.error(s.line+1, "Unexpected character: "+fmt.Sprintf("'%d'", c))
		} else {
			s.error(s.line+1, "Unexpected character: "+string(c))
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
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, "")
}

func (s *Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{tokenType, text, literal, s.line})
}

func (s *Scanner) scanString() {
	stringStartLine := s.line

	for !s.endOfString() && !s.isAtEnd() && s.peek() != '\n' {
		s.advance()
	}
	if s.isAtEnd() || s.peek() == '\n' {
		s.error(s.line+1, "Unterminated string at line "+fmt.Sprintf("%d", stringStartLine+1))
		return
	}
	s.advance()
	literal := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(STRING, s.replaceEscapeSequences(literal))
}

func (s *Scanner) scanComment() {
	for s.peek() != '\n' && !s.isAtEnd() {
		s.advance()
	}
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
		ffloat, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
		s.addTokenWithLiteral(NUMBER, ffloat)
	} else {
		iint, _ := strconv.Atoi(s.source[s.start:s.current])
		s.addTokenWithLiteral(NUMBER, iint)
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

	if tokenType, ok := keywords[text]; ok {
		s.addToken(tokenType)
	} else {
		s.addTokenWithLiteral(IDENTIFIER, text)
	}
}
