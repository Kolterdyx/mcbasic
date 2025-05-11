package symbol

import "github.com/Kolterdyx/mcbasic/internal/tokens"

type SymbolType string

const (
	_              SymbolType = ""
	FunctionSymbol            = "Function"
	StructSymbol              = "Struct"
	VariableSymbol            = "Variable"
	ImportSymbol              = "Import"
)

type Symbol struct {
	name             string
	stype            SymbolType
	declarationToken tokens.Token
	originFile       string
}

func NewSymbol(name string, stype SymbolType, declarationToken tokens.Token, originFile string) Symbol {
	return Symbol{
		name:             name,
		stype:            stype,
		declarationToken: declarationToken,
		originFile:       originFile,
	}
}

func (s Symbol) Name() string {
	return s.name
}

func (s Symbol) Type() SymbolType {
	return s.stype
}

func (s Symbol) DeclarationToken() tokens.Token {
	return s.declarationToken
}

func (s Symbol) OriginFile() string {
	return s.originFile
}
