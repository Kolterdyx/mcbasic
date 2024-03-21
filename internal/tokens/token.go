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
	Undefined TokenType = iota
	Identifier
	Number
	String
	Plus
	Minus
	Star
	Slash
	Percent
	Less
	LessEqual
	Greater
	GreaterEqual
	Equal
	EqualEqual
	Bang
	BangEqual
	ParenOpen
	ParenClose
	BraceOpen
	BraceClose
	Comma
	Semicolon

	Eof

	True
	False
	Print
	Let
	Def
	And
	Or
	Not
	If
	Else
	For
)

func (t TokenType) String() string {
	switch t {
	case Undefined:
		return "Undefined"
	case Identifier:
		return "Identifier"
	case Number:
		return "Number"
	case String:
		return "String"
	case Plus:
		return "Plus"
	case Minus:
		return "Minus"
	case Star:
		return "Star"
	case Slash:
		return "Slash"
	case Percent:
		return "Percent"
	case Less:
		return "Less"
	case LessEqual:
		return "LessEqual"
	case Greater:
		return "Greater"
	case GreaterEqual:
		return "GreaterEqual"
	case Equal:
		return "Equal"
	case EqualEqual:
		return "EqualEqual"
	case Bang:
		return "Bang"
	case BangEqual:
		return "BangEqual"
	case And:
		return "And"
	case Or:
		return "Or"
	case Not:
		return "Not"
	case If:
		return "If"
	case Else:
		return "Else"
	case For:
		return "For"
	case ParenOpen:
		return "ParenOpen"
	case ParenClose:
		return "ParenClose"
	case BraceOpen:
		return "BraceOpen"
	case BraceClose:
		return "BraceClose"
	case Comma:
		return "Comma"
	case Semicolon:
		return "Semicolon"
	case Eof:
		return "Eof"
	case Print:
		return "Print"
	case Let:
		return "Let"
	case Def:
		return "Def"
	case True:
		return "True"
	case False:
		return "False"
	default:
		return "Unknown"
	}
}

var Keywords = map[string]TokenType{
	"print": Print,
	"let":   Let,
	"def":   Def,
	"if":    If,
	"else":  Else,
	"for":   For,
	"and":   And,
	"or":    Or,
	"not":   Not,
	"true":  True,
	"false": False,
}
