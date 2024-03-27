package statements

type StmtType string

const (
	ExpressionStmtType          StmtType = "Expression"
	PrintStmtType               StmtType = "Print"
	ExecStmtType                StmtType = "Exec"
	VariableDeclarationStmtType StmtType = "VariableDeclaration"
	FunctionDeclarationStmtType StmtType = "FunctionDeclaration"
	VariableAssignmentStmtType  StmtType = "VariableAssignment"
	BlockStmtType               StmtType = "Block"
	WhileStmtType               StmtType = "While"
	IfStmtType                  StmtType = "If"
)

type StmtVisitor interface {
	VisitExpression(ExpressionStmt) interface{}
	VisitPrint(PrintStmt) interface{}
	VisitVariableDeclaration(VariableDeclarationStmt) interface{}
	VisitFunctionDeclaration(FunctionDeclarationStmt) interface{}
	VisitVariableAssignment(VariableAssignmentStmt) interface{}
	VisitBlock(BlockStmt) interface{}
	VisitWhile(WhileStmt) interface{}
	VisitExec(ExecStmt) interface{}
	VisitIf(IfStmt) interface{}
}

type Stmt interface {
	Accept(StmtVisitor) interface{}
	TType() StmtType
}
