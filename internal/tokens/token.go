package tokens

type TokenType int

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func (t Token) String() string {
	return t.Lexeme
}

const (
	UNDEFINED TokenType = iota
	IDENTIFIER
	NUMBER
	STRING
	PLUS
	MINUS
	STAR
	SLASH
	PERCENT
	LESS
	LESS_EQUAL
	GREATER
	GREATER_EQUAL
	EQUAL
	EQUAL_EQUAL
	BANG
	BANG_EQUAL
	PAREN_OPEN
	PAREN_CLOSE
	BRACE_OPEN
	BRACE_CLOSE
	COMMA
	SEMICOLON

	EOF

	// Keywords
	TRUE
	FALSE
	PRINT
	LET
	DEF
	AND
	OR
	NOT
	IF
	ELSE
	FOR
)

func (t TokenType) String() string {
	switch t {
	case UNDEFINED:
		return "UNDEFINED"
	case IDENTIFIER:
		return "IDENTIFIER"
	case NUMBER:
		return "NUMBER"
	case STRING:
		return "STRING"
	case PLUS:
		return "PLUS"
	case MINUS:
		return "MINUS"
	case STAR:
		return "STAR"
	case SLASH:
		return "SLASH"
	case PERCENT:
		return "PERCENT"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case NOT:
		return "NOT"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case FOR:
		return "FOR"
	case PAREN_OPEN:
		return "PAREN_OPEN"
	case PAREN_CLOSE:
		return "PAREN_CLOSE"
	case BRACE_OPEN:
		return "BRACE_OPEN"
	case BRACE_CLOSE:
		return "BRACE_CLOSE"
	case COMMA:
		return "COMMA"
	case SEMICOLON:
		return "SEMICOLON"
	case EOF:
		return "EOF"
	case PRINT:
		return "PRINT"
	case LET:
		return "LET"
	case DEF:
		return "DEF"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	default:
		return "UNKNOWN"
	}
}

var Keywords = map[string]TokenType{
	"print": PRINT,
	"let":   LET,
	"def":   DEF,
	"if":    IF,
	"else":  ELSE,
	"for":   FOR,
	"and":   AND,
	"or":    OR,
	"not":   NOT,
	"true":  TRUE,
	"false": FALSE,
}
