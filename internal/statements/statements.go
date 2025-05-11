package statements

import (
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
)

type StmtVisitor interface {
	VisitExpression(ExpressionStmt) interfaces.IRCode
	VisitVariableDeclaration(VariableDeclarationStmt) interfaces.IRCode
	VisitFunctionDeclaration(FunctionDeclarationStmt) interfaces.IRCode
	VisitVariableAssignment(VariableAssignmentStmt) interfaces.IRCode
	VisitStructDeclaration(StructDeclarationStmt) interfaces.IRCode
	VisitBlock(BlockStmt) interfaces.IRCode
	VisitWhile(WhileStmt) interfaces.IRCode
	VisitIf(IfStmt) interfaces.IRCode
	VisitReturn(ReturnStmt) interfaces.IRCode
	VisitScore(ScoreStmt) interfaces.IRCode
	VisitSetReturnFlag(SetReturnFlagStmt) interfaces.IRCode
	VisitImport(ImportStmt) interfaces.IRCode
}

type Stmt interface {
	ast.Node
	Accept(StmtVisitor) interfaces.IRCode
	ToString() string
}
