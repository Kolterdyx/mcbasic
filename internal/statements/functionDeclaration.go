package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
)

type FunctionDeclarationStmt struct {
	Stmt

	Name       tokens.Token
	Parameters []interfaces.FuncArg
	ReturnType types.ValueType
	Body       BlockStmt
}

func (f FunctionDeclarationStmt) Accept(visitor StmtVisitor) string {
	return visitor.VisitFunctionDeclaration(f)
}

func (f FunctionDeclarationStmt) StmtType() StmtType {
	return FunctionDeclarationStmtType
}

func (f FunctionDeclarationStmt) HasArg(arg string) bool {
	for _, p := range f.Parameters {
		if p.Name == arg {
			return true
		}
	}
	return false
}
