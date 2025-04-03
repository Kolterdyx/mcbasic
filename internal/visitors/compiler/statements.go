package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func (c *Compiler) VisitFunctionDeclaration(stmt statements.FunctionDeclarationStmt) interface{} {
	c.currentFunction = stmt
	c.currentScope = stmt.Name.Lexeme
	c.scope[stmt.Name.Lexeme] = []TypedIdentifier{}

	cmd := ""
	// For each parameter, copy the value to a variable with the same name
	for _, arg := range stmt.Parameters {
		cmd += c.opHandler.MoveConst(c.opHandler.Macro(arg.Name), ops.Cs(arg.Name))
	}

	var source = cmd + stmt.Body.Accept(c).(string)

	args := make([]interfaces.FuncArg, 0)
	for _, arg := range stmt.Parameters {
		args = append(args, interfaces.FuncArg{Name: arg.Name, Type: arg.Type})
		c.scope[stmt.Name.Lexeme] = append(c.scope[stmt.Name.Lexeme],
			TypedIdentifier{
				arg.Name,
				arg.Type,
			})
	}
	// add function to all scopes
	for scope := range c.scope {
		c.scope[scope] = append(c.scope[scope],
			TypedIdentifier{
				stmt.Name.Lexeme,
				stmt.ReturnType,
			})
	}

	// function wrapper that automatically loads the __call__ parameter
	wrapperSource := ""
	for _, arg := range stmt.Parameters {
		wrapperSource += c.opHandler.LoadArgConst("internal/"+stmt.Name.Lexeme+"__wrapped", arg.Name, c.opHandler.Macro(arg.Name))
	}
	wrapperSource += c.opHandler.Call("internal/"+stmt.Name.Lexeme+"__wrapped", "")
	c.createFunction(stmt.Name.Lexeme, c.opHandler.MacroReplace(wrapperSource), args, stmt.ReturnType)
	c.createFunction("internal/"+stmt.Name.Lexeme+"__wrapped", c.opHandler.MacroReplace(source), args, stmt.ReturnType)
	return ""
}

func (c *Compiler) VisitBlock(stmt statements.BlockStmt) interface{} {
	section := ""
	for _, s := range stmt.Statements {
		section += s.Accept(c).(string)
	}
	//c.currentFunctionSections = append(c.currentFunctionSections, section)
	return section // c.opHandler.CallSection(c.currentFunction.Name.Lexeme, len(c.currentFunctionSections)-1)
}

func (c *Compiler) VisitExpression(stmt statements.ExpressionStmt) interface{} {
	return stmt.Expression.Accept(c).(string)
}

func (c *Compiler) VisitReturn(stmt statements.ReturnStmt) interface{} {
	cmd := ""
	if stmt.Expression.ReturnType() != c.currentFunction.ReturnType {
		log.Fatalln("Return type does not match function return type")
	}
	cmd += stmt.Expression.Accept(c).(string)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.RET)
	cmd += c.opHandler.Return()
	return cmd
}

func (c *Compiler) VisitVariableDeclaration(stmt statements.VariableDeclarationStmt) interface{} {
	cmd := ""
	if stmt.Initializer != nil {
		cmd += stmt.Initializer.Accept(c).(string)
	}
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(stmt.Name.Lexeme))
	c.scope[c.currentScope] = append(c.scope[c.currentScope],
		TypedIdentifier{
			stmt.Name.Lexeme,
			stmt.Type,
		})
	return cmd
}

func (c *Compiler) VisitVariableAssignment(stmt statements.VariableAssignmentStmt) interface{} {
	cmd := ""

	fmt.Println(stmt.Name.Lexeme, stmt.Value.ReturnType(), c.getReturnType(stmt.Name.Lexeme))
	if stmt.Value.ReturnType() != c.getReturnType(stmt.Name.Lexeme) {
		log.Fatalln("Assignment type mismatch")
	}
	cmd += stmt.Value.Accept(c).(string)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(stmt.Name.Lexeme))
	return cmd
}

func (c *Compiler) VisitIf(stmt statements.IfStmt) interface{} {
	cmd := ""
	cmd += stmt.Condition.Accept(c).(string)
	thenSource := stmt.ThenBranch.Accept(c).(string)
	elseSource := ""
	if reflect.ValueOf(stmt.ElseBranch) != reflect.ValueOf(statements.BlockStmt{}) {
		elseSource = stmt.ElseBranch.Accept(c).(string)
	}
	condVar := c.newRegister(ops.RX)
	cmd += c.opHandler.MoveScore(ops.Cs(ops.RX), ops.Cs(condVar))
	cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 1", ops.Cs(condVar), c.Namespace), true, thenSource)
	cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 1", ops.Cs(condVar), c.Namespace), false, elseSource)

	return cmd
}
