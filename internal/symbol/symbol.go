package symbol

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type Type string

const (
	_              Type = ""
	FunctionSymbol      = "Function"
	StructSymbol        = "Struct"
	VariableSymbol      = "Variable"
)

type Symbol struct {
	name            string
	stype           Type
	declarationNode ast.Node
	originFile      string
}

func NewSymbol(name string, stype Type, declarationNode ast.Node, originFile string) Symbol {
	return Symbol{
		name:            name,
		stype:           stype,
		declarationNode: declarationNode,
		originFile:      originFile,
	}
}

func (s Symbol) Name() string {
	return s.name
}

func (s Symbol) Type() Type {
	return s.stype
}

func (s Symbol) OriginFile() string {
	return s.originFile
}

func (s Symbol) DeclarationNode() ast.Node {
	return s.declarationNode
}

func (s Symbol) ValueType() types.ValueType {
	switch s.stype {
	case VariableSymbol:
		if variable, ok := s.declarationNode.(statements.VariableDeclarationStmt); ok {
			return variable.ValueType
		}
	case FunctionSymbol:
		if function, ok := s.declarationNode.(statements.FunctionDeclarationStmt); ok {
			return function.ReturnType
		}
	case StructSymbol:
		if structDecl, ok := s.declarationNode.(statements.StructDeclarationStmt); ok {
			return structDecl.StructType
		}
	}
	return nil
}
