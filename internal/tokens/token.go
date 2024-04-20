package tokens

type TokenType int

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal string
	Line    int
	Column  int
}

func (t Token) String() string {
	return t.Lexeme
}

const (
	Undefined TokenType = iota
	Identifier
	Int
	Fixed
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
	Colon
	BracketOpen
	BracketClose

	Eof

	True
	False
	Let
	Def
	And
	Or
	If
	Else
	While
	Return

	IntType
	FixedType
	StringType
)

func (t TokenType) String() string {
	switch t {
	case Undefined:
		return "Undefined"
	case Identifier:
		return "Identifier"
	case Int:
		return "Int"
	case Fixed:
		return "Fixed"
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
		return "Eq"
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
	case If:
		return "If"
	case Else:
		return "Else"
	case While:
		return "While"
	case Return:
		return "Return"
	case ParenOpen:
		return "ParenOpen"
	case ParenClose:
		return "ParenClose"
	case BraceOpen:
		return "BraceOpen"
	case BraceClose:
		return "BraceClose"
	case BracketOpen:
		return "BracketOpen"
	case BracketClose:
		return "BracketClose"
	case Comma:
		return "Comma"
	case Semicolon:
		return "Semicolon"
	case Colon:
		return "Colon"
	case Eof:
		return "Eof"
	case Let:
		return "Let"
	case Def:
		return "Def"
	case True:
		return "True"
	case False:
		return "False"
	case IntType:
		return "IntType"
	case FixedType:
		return "FixedType"
	case StringType:
		return "StringType"
	default:
		return "Unknown"
	}
}

var Keywords = map[string]TokenType{
	"let":    Let,
	"func":   Def,
	"if":     If,
	"else":   Else,
	"and":    And,
	"or":     Or,
	"true":   True,
	"false":  False,
	"while":  While,
	"int":    IntType,
	"fixed":  FixedType,
	"str":    StringType,
	"return": Return,
}
