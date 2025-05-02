package il

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"strings"
)

func (c *Compiler) VisitExpression(stmt statements.ExpressionStmt) (cmd string) {
	return stmt.Expression.Accept(c)
}

func (c *Compiler) VisitVariableDeclaration(stmt statements.VariableDeclarationStmt) (cmd string) {
	if stmt.Initializer != nil {
		cmd += stmt.Initializer.Accept(c)
		cmd += c.CopyVar(RX, stmt.Name.Lexeme)
	} else {
		cmd = c.SetVar(stmt.Name.Lexeme, stmt.Type.ToNBT())
	}
	return
}

func (c *Compiler) VisitFunctionDeclaration(stmt statements.FunctionDeclarationStmt) (cmd string) {

	c.currentScope = stmt.Name.Lexeme
	c.branchCounter = 0
	c.scopes[c.currentScope] = make([]interfaces.TypedIdentifier, 0)

	cmd += c.Func(c.currentScope)

	// For each parameter, copy the value to a variable with the same name and add it to the scope
	for _, arg := range stmt.Parameters {
		macro := fmt.Sprintf("$(%s)", arg.Name)
		if arg.Type.Equals(types.StringType) {
			cmd += c.Set(c.varPath(arg.Name), nbt.NewString(macro))
		} else {
			cmd += c.Set(c.varPath(arg.Name), nbt.NewAny(macro))
		}
		c.scopes[c.currentScope] = append(c.scopes[c.currentScope],
			interfaces.TypedIdentifier{
				Name: arg.Name,
				Type: arg.Type,
			})
	}

	// add function to all scopes
	for scope := range c.scopes {
		c.scopes[scope] = append(c.scopes[scope],
			interfaces.TypedIdentifier{
				Name: stmt.Name.Lexeme,
				Type: stmt.ReturnType,
			})
	}

	cmd += stmt.Body.Accept(c)
	cmd += c.Ret()
	splitSrc := strings.Split(cmd, "\n")
	for i, line := range splitSrc {
		if !strings.HasPrefix(line, Func) {
			splitSrc[i] = "\t" + line
		}
	}
	cmd = strings.TrimSpace(strings.Join(splitSrc, "\n")) + "\n"
	return
}

func (c *Compiler) VisitVariableAssignment(stmt statements.VariableAssignmentStmt) (cmd string) {
	isIndexedAssignment := len(stmt.Accessors) > 0
	if !stmt.Value.ReturnType().Equals(c.getReturnType(stmt.Name.Lexeme)) && !isIndexedAssignment {
		c.error(stmt.Name.SourceLocation, fmt.Sprintf("Assignment type mismatch: %v != %v", c.getReturnType(stmt.Name.Lexeme).ToString(), stmt.Value.ReturnType().ToString()))
	}
	cmd += stmt.Value.Accept(c)
	valueReg := c.makeReg(RX)
	cmd += c.CopyVar(RX, valueReg)
	if isIndexedAssignment {
		pathReg := c.makeReg(RX)
		cmd += c.Set(pathReg, nbt.NewString(""))
		for _, accessor := range stmt.Accessors {
			switch accessor := accessor.(type) {
			case statements.IndexAccessor:
				indexReg := c.makeReg(RX)
				cmd += accessor.Index.Accept(c)
				cmd += c.CopyVar(RX, indexReg)
				// wrap index in brackets and append to pathReg
				cmd += c.MakeIndex(indexReg, indexReg)
				cmd += c.StringConcat(pathReg, indexReg, pathReg)
			case statements.FieldAccessor:
				fieldReg := c.makeReg(RX)
				cmd += c.Set(fieldReg, nbt.NewString(accessor.ToString()))
				cmd += c.StringConcat(pathReg, fieldReg, pathReg)
			}
		}
		cmd += c.PathSet(stmt.Name.Lexeme, pathReg, valueReg)
	} else {
		cmd += c.CopyVar(valueReg, stmt.Name.Lexeme)
	}
	return
}

func (c *Compiler) VisitStructDeclaration(stmt statements.StructDeclarationStmt) (cmd string) {
	return c.Set(c.structPath(stmt.Name.Lexeme), stmt.StructType.ToNBT())
}

func (c *Compiler) VisitBlock(stmt statements.BlockStmt) (cmd string) {
	for _, statement := range stmt.Statements {
		cmd += statement.Accept(c)
	}
	return
}

func (c *Compiler) VisitWhile(stmt statements.WhileStmt) (cmd string) {
	//TODO implement me
	panic("implement me")
}

func (c *Compiler) VisitIf(stmt statements.IfStmt) (cmd string) {

	c.branchCounter++
	thenBranchName := fmt.Sprintf("%s_%d", c.currentScope, c.branchCounter)
	c.branchCounter++
	elseBranchName := fmt.Sprintf("%s_%d", c.currentScope, c.branchCounter)

	cmd += c.Score(RETF, nbt.NewInt(0))
	c.compiledFunctions[thenBranchName] = c.makeBranchFunction(thenBranchName, stmt.ThenBranch).Accept(c)

	hasElseBranch := stmt.ElseBranch != nil
	if hasElseBranch {
		c.compiledFunctions[elseBranchName] = c.makeBranchFunction(thenBranchName, *stmt.ElseBranch).Accept(c)
	}
	cmd += stmt.Condition.Accept(c)
	condVar := c.makeReg(RX)
	cmd += c.CopyVar(RX, condVar)
	cmd += c.If(condVar, c.Call(thenBranchName))
	cmd += c.If(RETF, c.Ret())
	if hasElseBranch {
		cmd += c.Unless(condVar, c.Call(elseBranchName))
		cmd += c.If(RETF, c.Ret())
	}
	return
}

func (c *Compiler) VisitReturn(stmt statements.ReturnStmt) (cmd string) {
	if stmt.Expression != nil {
		cmd += stmt.Expression.Accept(c)
		cmd += c.CopyVar(RX, RET)
	}
	cmd += c.Ret()
	return
}

func (c *Compiler) VisitScore(stmt statements.ScoreStmt) (cmd string) {
	cmd += c.Score(stmt.Target, nbt.NewInt(stmt.Score))
	return
}
