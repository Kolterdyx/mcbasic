package statements

import "github.com/Kolterdyx/mcbasic/internal/interfaces"

type StmtType string

const (
	ExpressionStmtType          StmtType = "Expression"
	VariableDeclarationStmtType StmtType = "VariableDeclaration"
	FunctionDeclarationStmtType StmtType = "FunctionDeclaration"
	StructDeclarationStmtType   StmtType = "StructDeclaration"
	VariableAssignmentStmtType  StmtType = "VariableAssignment"
	BlockStmtType               StmtType = "Block"
	WhileStmtType               StmtType = "While"
	IfStmtType                  StmtType = "If"
	ReturnStmtType              StmtType = "Return"
	ScoreStmtType               StmtType = "Score"
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
}

type Stmt interface {
	Accept(StmtVisitor) interfaces.IRCode
	StmtType() StmtType
	ToString() string
}
