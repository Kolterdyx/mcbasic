package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func (c *Compiler) VisitFunctionDeclaration(stmt statements.FunctionDeclarationStmt) string {
	c.currentFunction = stmt
	c.currentScope = stmt.Name.Lexeme
	c.scope[stmt.Name.Lexeme] = []TypedIdentifier{}

	cmd := ""
	// For each parameter, copy the value to a variable with the same name
	for _, arg := range stmt.Parameters {
		cmd += c.opHandler.MakeConst(c.opHandler.Macro(arg.Name), ops.Cs(arg.Name))
	}

	var source = cmd + stmt.Body.Accept(c)

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

func (c *Compiler) VisitBlock(stmt statements.BlockStmt) string {
	section := ""
	for _, s := range stmt.Statements {
		section += s.Accept(c)
	}
	//c.currentFunctionSections = append(c.currentFunctionSections, section)
	return section // c.opHandler.CallSection(c.currentFunction.Name.Lexeme, len(c.currentFunctionSections)-1)
}

func (c *Compiler) VisitExpression(stmt statements.ExpressionStmt) string {
	return stmt.Expression.Accept(c)
}

func (c *Compiler) VisitReturn(stmt statements.ReturnStmt) string {
	cmd := ""
	if stmt.Expression.ReturnType() != c.currentFunction.ReturnType {
		log.Fatalln("Return type does not match function return type")
	}
	cmd += stmt.Expression.Accept(c)
	cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.RET)
	cmd += c.opHandler.Return()
	return cmd
}

func (c *Compiler) VisitVariableDeclaration(stmt statements.VariableDeclarationStmt) string {
	cmd := ""
	if stmt.Initializer != nil {
		cmd += stmt.Initializer.Accept(c)
		cmd += c.opHandler.Move(ops.Cs(ops.RX), ops.Cs(stmt.Name.Lexeme))
	} else {
		switch stmt.Type {
		case types.IntType:
			cmd += c.opHandler.MakeConst("0L", ops.Cs(stmt.Name.Lexeme), false)
		case types.DoubleType:
			cmd += c.opHandler.MakeConst("0.0d", ops.Cs(stmt.Name.Lexeme), false)
		case types.StringType:
			cmd += c.opHandler.MakeConst("\"\"", ops.Cs(stmt.Name.Lexeme))
		default:
			if reflect.TypeOf(stmt.Type) == reflect.TypeOf(types.ListTypeStruct{}) {
				cmd += c.opHandler.MakeList(ops.Cs(stmt.Name.Lexeme))
			} else if reflect.TypeOf(stmt.Type) == reflect.TypeOf(types.StructTypeStruct{}) {
				cmd += c.opHandler.MoveRaw(
					fmt.Sprintf("%s:data", c.Namespace),
					fmt.Sprintf("%s.%s", ops.StructPath, stmt.Type),
					fmt.Sprintf("%s:data", c.Namespace),
					fmt.Sprintf("%s.%s", ops.VarPath, ops.Cs(stmt.Name.Lexeme)),
				)
			} else {
				c.error(stmt.Name.SourceLocation, fmt.Sprintf("Invalid type: %v", stmt.Type))
			}
		}
	}
	c.scope[c.currentScope] = append(c.scope[c.currentScope],
		TypedIdentifier{
			stmt.Name.Lexeme,
			stmt.Type,
		})
	return cmd
}

func (c *Compiler) VisitVariableAssignment(stmt statements.VariableAssignmentStmt) string {
	cmd := ""
	isIndexedAssignment := stmt.Index != nil
	if stmt.Value.ReturnType() != c.getReturnType(stmt.Name.Lexeme) && !isIndexedAssignment {
		c.error(stmt.Name.SourceLocation, fmt.Sprintf("Assignment type mismatch: %v != %v", c.getReturnType(stmt.Name.Lexeme), stmt.Value.ReturnType()))
	}
	cmd += stmt.Value.Accept(c)
	valueReg := ops.Cs(c.newRegister(ops.RX))
	cmd += c.opHandler.Move(ops.Cs(ops.RX), valueReg)
	if isIndexedAssignment {
		cmd += stmt.Index.Accept(c)
		indexReg := ops.Cs(c.newRegister(ops.RX))
		cmd += c.opHandler.Move(ops.Cs(ops.RX), indexReg)
		cmd += c.opHandler.SetListIndex(ops.Cs(stmt.Name.Lexeme), indexReg, valueReg)
	} else {
		cmd += c.opHandler.Move(valueReg, ops.Cs(stmt.Name.Lexeme))
	}
	return cmd
}

func (c *Compiler) VisitIf(stmt statements.IfStmt) string {
	cmd := ""
	cmd += stmt.Condition.Accept(c)
	thenSource := stmt.ThenBranch.Accept(c)
	elseSource := ""
	if reflect.ValueOf(stmt.ElseBranch) != reflect.ValueOf(statements.BlockStmt{}) {
		elseSource = stmt.ElseBranch.Accept(c)
	}
	condVar := c.newRegister(ops.RX)
	cmd += c.opHandler.MoveScore(ops.Cs(ops.RX), ops.Cs(condVar))
	cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 1", ops.Cs(condVar), c.Namespace), true, thenSource)
	cmd += c.opHandler.ExecCond(fmt.Sprintf("score %s %s matches 1", ops.Cs(condVar), c.Namespace), false, elseSource)

	return cmd
}

func (c *Compiler) VisitStructDeclaration(stmt statements.StructDeclarationStmt) string {
	return c.opHandler.StructDefine(c.getReturnType(stmt.Name.Lexeme).(types.StructTypeStruct))
}
