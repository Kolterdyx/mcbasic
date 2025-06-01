package tokens

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type Token struct {
	interfaces.SourceLocation
	Type    interfaces.TokenType
	Lexeme  string
	Literal string
}

func (t Token) String() string {
	return t.Lexeme
}

const (
	Undefined interfaces.TokenType = iota
	Identifier
	Int
	Double
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
	Dot
	Comma
	Semicolon
	Colon
	BracketOpen
	BracketClose

	Eof

	True
	False
	Let
	Func
	And
	Or
	If
	Else
	While
	Return
	Import
	Exec

	IntType
	DoubleType
	StringType
	VoidType

	Struct
)

var Keywords = map[string]interfaces.TokenType{
	"let":    Let,
	"func":   Func,
	"if":     If,
	"else":   Else,
	"and":    And,
	"or":     Or,
	"true":   True,
	"false":  False,
	"while":  While,
	"int":    IntType,
	"double": DoubleType,
	"str":    StringType,
	"void":   VoidType,
	"return": Return,
	"struct": Struct,
	"import": Import,
	"exec":   Exec,
}
