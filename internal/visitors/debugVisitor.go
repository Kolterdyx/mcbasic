package visitors

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

func (d DebugVisitor) VisitExpression(e statements.ExpressionStmt) interface{} {
	return e.Expression.Accept(d)
}

func (d DebugVisitor) VisitPrint(p statements.PrintStmt) interface{} {
	return p.Expression.Accept(d)
}

func (d DebugVisitor) VisitVariable(v expressions.VariableExpr) interface{} {
	return v.Name.Lexeme
}

func (d DebugVisitor) VisitVariableDeclaration(v statements.VariableDeclarationStmt) interface{} {
	res := "let " + v.Name.Lexeme
	if v.Initializer != nil {
		res += " = " + v.Initializer.Accept(d).(string)
	}
	return res + ";"
}

func (d DebugVisitor) VisitVariableAssignment(v statements.VariableAssignmentStmt) interface{} {
	return v.Name.Lexeme + " = " + v.Value.Accept(d).(string) + ";"
}

func (d DebugVisitor) VisitFunctionDeclaration(f statements.FunctionDeclarationStmt) interface{} {
	res := fmt.Sprintf("def %s(", f.Name.Lexeme)
	for i, p := range f.Parameters {
		res += p.Lexeme
		if i < len(f.Parameters)-1 {
			res += ", "
		}
	}
	res += ") " + f.Body.Accept(d).(string)
	return res
}

func (d DebugVisitor) VisitFunctionCall(f expressions.FunctionCallExpr) interface{} {
	res := fmt.Sprintf("%s(", f.Callee.Lexeme)
	for i, a := range f.Arguments {
		res += a.Accept(d).(string)
		if i < len(f.Arguments)-1 {
			res += ", "
		}
	}
	res += ")"
	return res
}

func (d DebugVisitor) VisitBlock(b statements.BlockStmt) interface{} {
	res := "{\n"
	for _, s := range b.Statements {
		res += s.Accept(d).(string) + "\n"
	}
	res += "}"
	return res
}
