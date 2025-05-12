package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
)

type StructDeclarationStmt struct {
	Stmt

	Name tokens.Token
}

func (s StructDeclarationStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitStructDeclaration(s)
}

func (s StructDeclarationStmt) Type() ast.NodeType {
	return ast.StructDeclarationStatement
}

func (s StructDeclarationStmt) ToString() string {
	return "struct " + s.Name.Lexeme
}
