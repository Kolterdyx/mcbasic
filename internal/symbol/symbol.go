package symbol

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"path"
)

type Type string

const (
	_              Type = ""
	FunctionSymbol      = "Function"
	StructSymbol        = "Struct"
	VariableSymbol      = "Variable"
	LiteralSymbol       = "Literal"
	ImportSymbol        = "Import"
)

type Symbol struct {
	name            string
	stype           Type
	declarationNode ast.Node
	valueType       types.ValueType
}

func NewSymbol(name string, stype Type, declarationNode ast.Node, valueType types.ValueType) Symbol {
	return Symbol{
		name:            name,
		stype:           stype,
		declarationNode: declarationNode,
		valueType:       valueType,
	}
}

func (s Symbol) Name() string {
	return s.name
}

func (s Symbol) Type() Type {
	return s.stype
}

func (s Symbol) DeclarationNode() ast.Node {
	return s.declarationNode
}

func (s Symbol) ValueType() types.ValueType {
	return s.valueType
}

func (s Symbol) AsImportedFrom(importName string) Symbol {
	return NewSymbol(
		fmt.Sprintf("%s.%s", path.Base(importName), s.name),
		s.stype,
		s.declarationNode,
		s.valueType,
	)
}
