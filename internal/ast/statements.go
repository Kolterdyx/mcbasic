package ast

import "fmt"

type StatementVisitor interface {
	VisitExpression(ExpressionStmt) any
	VisitVariableDeclaration(VariableDeclarationStmt) any
	VisitFunctionDeclaration(FunctionDeclarationStmt) any
	VisitVariableAssignment(VariableAssignmentStmt) any
	VisitStructDeclaration(StructDeclarationStmt) any
	VisitBlock(BlockStmt) any
	VisitWhile(WhileStmt) any
	VisitIf(IfStmt) any
	VisitReturn(ReturnStmt) any
	VisitSetReturnFlag(SetReturnFlagStmt) any
	VisitImport(ImportStmt) any
	VisitExec(ExecStmt) any
}

type Statement interface {
	Node
	Accept(StatementVisitor) any
	ToString() string
}

func AcceptStmt[T any](stmt Statement, v StatementVisitor) T {
	res := stmt.Accept(v)
	if res == nil {
		var zero T
		return zero
	}
	val, ok := res.(T)
	if !ok {
		panic(fmt.Sprintf("unexpected type: got %T, want %T", res, val))
	}
	return val
}

type Source []Statement
