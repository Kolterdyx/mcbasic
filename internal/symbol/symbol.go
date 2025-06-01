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
	alias           string
	stype           Type
	declarationNode ast.Node
	valueType       types.ValueType
}

func NewSymbol(name string, stype Type, declarationNode ast.Node, valueType types.ValueType) Symbol {
	return Symbol{
		name:            name,
		alias:           name,
		stype:           stype,
		declarationNode: declarationNode,
		valueType:       valueType,
	}
}

func NewAlias(alias string, sym Symbol) Symbol {
	return Symbol{
		name:            sym.name,
		alias:           alias,
		stype:           sym.stype,
		declarationNode: sym.declarationNode,
		valueType:       sym.valueType,
	}
}

func (s Symbol) Alias() string {
	return s.alias
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
