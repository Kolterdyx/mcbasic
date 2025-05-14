package ast

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type VariableDeclarationStmt struct {
	Name        tokens.Token
	Initializer Expr
	ValueType   types.ValueType
}

func (v VariableDeclarationStmt) Accept(visitor StatementVisitor) any {
	return visitor.VisitVariableDeclaration(v)
}

func (v VariableDeclarationStmt) Type() NodeType {
	return VariableDeclarationStatement
}

func (v VariableDeclarationStmt) ToString() string {
	if v.Initializer != nil {
		return v.Name.Lexeme + " = " + v.Initializer.ToString() + ";"
	}
	return v.Name.Lexeme + ";"
}

func (v VariableDeclarationStmt) GetSourceLocation() interfaces.SourceLocation {
	return v.Name.SourceLocation
}
