package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type VariableDeclarationStmt struct {
	Stmt

	Name        tokens.Token
	Initializer expressions.Expr
	ValueType   types.ValueType
}

func (v VariableDeclarationStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitVariableDeclaration(v)
}

func (v VariableDeclarationStmt) Type() ast.NodeType {
	return ast.VariableDeclarationStatement
}

func (v VariableDeclarationStmt) ToString() string {
	if v.Initializer != nil {
		return v.Name.Lexeme + " = " + v.Initializer.ToString() + ";"
	}
	return v.Name.Lexeme + ";"
}
