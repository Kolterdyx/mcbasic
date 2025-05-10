package statements

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
	VisitExpression(ExpressionStmt) (cmd string)
	VisitVariableDeclaration(VariableDeclarationStmt) (cmd string)
	VisitFunctionDeclaration(FunctionDeclarationStmt) (cmd string)
	VisitVariableAssignment(VariableAssignmentStmt) (cmd string)
	VisitStructDeclaration(StructDeclarationStmt) (cmd string)
	VisitBlock(BlockStmt) (cmd string)
	VisitWhile(WhileStmt) (cmd string)
	VisitIf(IfStmt) (cmd string)
	VisitReturn(ReturnStmt) (cmd string)
	VisitScore(ScoreStmt) (cmd string)
}

type Stmt interface {
	Accept(StmtVisitor) (cmd string)
	StmtType() StmtType
}
