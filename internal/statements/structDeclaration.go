package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type StructDeclarationStmt struct {
	Stmt

	Name       tokens.Token
	StructType types.StructTypeStruct
}

func (s StructDeclarationStmt) Accept(visitor StmtVisitor) interfaces.IRCode {
	return visitor.VisitStructDeclaration(s)
}

func (s StructDeclarationStmt) Type() ast.NodeType {
	return ast.StructDeclarationStatement
}

func (s StructDeclarationStmt) ToString() string {
	body := ""
	for _, fieldName := range s.StructType.GetFieldNames() {
		field, _ := s.StructType.GetField(fieldName)
		body += fieldName + " " + field.ToString() + ";\n"
	}
	if len(body) > 0 {
		body = body[:len(body)-2]
	}
	return "struct " + s.Name.Lexeme + " {\n" + body + "\n}"
}
