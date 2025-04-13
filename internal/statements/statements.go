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
)

type StmtVisitor interface {
	VisitExpression(ExpressionStmt) string
	VisitVariableDeclaration(VariableDeclarationStmt) string
	VisitFunctionDeclaration(FunctionDeclarationStmt) string
	VisitVariableAssignment(VariableAssignmentStmt) string
	VisitStructDeclaration(StructDeclarationStmt) string
	VisitBlock(BlockStmt) string
	VisitWhile(WhileStmt) string
	VisitIf(IfStmt) string
	VisitReturn(ReturnStmt) string
}

type Stmt interface {
	Accept(StmtVisitor) string
	TType() StmtType
}
