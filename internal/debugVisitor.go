package internal

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/statements"
)

type DebugVisitor struct {
	expressions.ExprVisitor
	statements.StmtVisitor
}

func (d DebugVisitor) VisitBinary(b expressions.BinaryExpr) interface{} {
	return fmt.Sprintf("%v %s %v", b.Left.Accept(d), b.Operator.Lexeme, b.Right.Accept(d))
}

func (d DebugVisitor) VisitGrouping(g expressions.GroupingExpr) interface{} {
	return fmt.Sprintf("( %v )", g.Expression.Accept(d))
}

func (d DebugVisitor) VisitLiteral(l expressions.LiteralExpr) interface{} {
	return fmt.Sprintf("%v", l.Value)
}

func (d DebugVisitor) VisitUnary(u expressions.UnaryExpr) interface{} {
	return fmt.Sprintf("%s%v", u.Operator.Lexeme, u.Expression.Accept(d))
}

func (d DebugVisitor) VisitVariable(v expressions.VariableExpr) interface{} {
	return fmt.Sprintf("%s", v.Name.Lexeme)
}

func (d DebugVisitor) VisitFunctionDeclaration(f statements.FunctionDeclarationStmt) {
	res := fmt.Sprintf("def %s(", f.Name.Lexeme)
	for i, p := range f.Parameters {
		res += p.Lexeme
		if i < len(f.Parameters)-1 {
			res += ", "
		}
	}
	res += ") {"
	fmt.Println(res)
	for _, s := range f.Body {
		s.Accept(d)
	}
	fmt.Println("}")
}

func (d DebugVisitor) VisitVariableDeclaration(v statements.VariableDeclarationStmt) {
	if v.Initializer == nil {
		fmt.Println("let " + v.Name.Lexeme + ";")
		return
	}
	fmt.Println("let " + v.Name.Lexeme + " = " + v.Initializer.Accept(d).(string) + ";")
}

func (d DebugVisitor) VisitExpression(e statements.ExpressionStmt) {
	fmt.Println(e.Expression.Accept(d))
}

func (d DebugVisitor) VisitPrint(p statements.PrintStmt) {
	fmt.Println(p.Expression.Accept(d))
}
