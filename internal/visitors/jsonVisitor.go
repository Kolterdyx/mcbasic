package visitors

import (
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/statements"
)

type JsonVisitor struct {
	expressions.ExprVisitor
	statements.StmtVisitor
}

func (j JsonVisitor) VisitBinary(b expressions.BinaryExpr) interface{} {
	return map[string]interface{}{
		"type":     b.Type(),
		"left":     b.Left.Accept(j),
		"operator": b.Operator.Lexeme,
		"right":    b.Right.Accept(j),
	}
}

func (j JsonVisitor) VisitGrouping(g expressions.GroupingExpr) interface{} {
	return map[string]interface{}{
		"type":       g.Type(),
		"expression": g.Expression.Accept(j),
	}
}

func (j JsonVisitor) VisitLiteral(l expressions.LiteralExpr) interface{} {
	return map[string]interface{}{
		"type":  l.Type(),
		"value": l.Value,
	}
}

func (j JsonVisitor) VisitUnary(u expressions.UnaryExpr) interface{} {
	return map[string]interface{}{
		"type":       u.Type(),
		"operator":   u.Operator.Lexeme,
		"expression": u.Expression.Accept(j),
	}
}

func (j JsonVisitor) VisitFunctionCall(f expressions.FunctionCallExpr) interface{} {
	args := make([]interface{}, 0)
	for _, a := range f.Arguments {
		args = append(args, a.Accept(j))
	}
	return map[string]interface{}{
		"type":      f.Type(),
		"callee":    f.Callee.Lexeme,
		"arguments": args,
	}
}

func (j JsonVisitor) VisitExpression(e statements.ExpressionStmt) interface{} {
	return map[string]interface{}{
		"type":       e.Type(),
		"expression": e.Expression.Accept(j),
	}
}

func (j JsonVisitor) VisitPrint(p statements.PrintStmt) interface{} {
	return map[string]interface{}{
		"type":       p.Type(),
		"expression": p.Expression.Accept(j),
	}
}

func (j JsonVisitor) VisitVariable(v expressions.VariableExpr) interface{} {
	return map[string]interface{}{
		"type": v.Type(),
		"name": v.Name.Lexeme,
	}
}

func (j JsonVisitor) VisitVariableDeclaration(v statements.VariableDeclarationStmt) interface{} {
	if v.Initializer == nil {
		return map[string]interface{}{
			"type": v.Type(),
			"name": v.Name.Lexeme,
		}
	}
	return map[string]interface{}{
		"type":        v.Type(),
		"name":        v.Name.Lexeme,
		"initializer": v.Initializer.Accept(j),
	}
}

func (j JsonVisitor) VisitVariableAssignment(v statements.VariableAssignmentStmt) interface{} {
	return map[string]interface{}{
		"type":  v.Type(),
		"name":  v.Name.Lexeme,
		"value": v.Value.Accept(j),
	}
}

func (j JsonVisitor) VisitFunctionDeclaration(f statements.FunctionDeclarationStmt) interface{} {
	return map[string]interface{}{
		"type":       f.Type(),
		"name":       f.Name.Lexeme,
		"parameters": f.Parameters,
		"body":       f.Body.Accept(j),
	}
}

func (j JsonVisitor) VisitBlock(b statements.BlockStmt) interface{} {
	stmts := make([]interface{}, 0)
	for _, s := range b.Statements {
		stmts = append(stmts, s.Accept(j))
	}
	return map[string]interface{}{
		"type":       b.Type(),
		"statements": stmts,
	}
}

func (j JsonVisitor) VisitWhile(w statements.WhileStmt) interface{} {
	return map[string]interface{}{
		"type":      w.Type(),
		"condition": w.Condition.Accept(j),
		"body":      w.Body.Accept(j),
	}
}
