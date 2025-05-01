package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/compiler/mapping"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/nbt"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
	"reflect"
)

func (c *Compiler) VisitFunctionDeclaration(stmt statements.FunctionDeclarationStmt) string {
	c.currentFunction = stmt
	c.currentScope = stmt.Name.Lexeme
	c.scopes[stmt.Name.Lexeme] = []TypedIdentifier{}

	cmd := ""
	// For each parameter, copy the value to a variable with the same name
	for _, arg := range stmt.Parameters {
		var value nbt.Value = nbt.NewAny(c.commandMapper.MacroWrapper(arg.Name))
		if arg.Type == types.StringType {
			value = nbt.NewString(c.commandMapper.MacroWrapper(arg.Name))
		}
		cmd += c.commandMapper.MakeConst(value, c.commandMapper.Cs(arg.Name))
	}

	var source = cmd + stmt.Body.Accept(c)
	source += c.commandMapper.Return()

	args := make([]interfaces.FuncArg, 0)
	for _, arg := range stmt.Parameters {
		args = append(args, interfaces.FuncArg{Name: arg.Name, Type: arg.Type})
		c.scopes[stmt.Name.Lexeme] = append(c.scopes[stmt.Name.Lexeme],
			TypedIdentifier{
				arg.Name,
				arg.Type,
			})
	}
	// add function to all scopes
	for scope := range c.scopes {
		c.scopes[scope] = append(c.scopes[scope],
			TypedIdentifier{
				stmt.Name.Lexeme,
				stmt.ReturnType,
			})
	}

	// function wrapper that automatically loads the __call__ parameter
	wrapperSource := ""
	for _, arg := range stmt.Parameters {
		var name nbt.Value = nbt.NewAny(c.commandMapper.MacroWrapper(arg.Name))
		if arg.Type == types.StringType {
			name = nbt.NewString(c.commandMapper.MacroWrapper(arg.Name))
		}
		wrapperSource += c.commandMapper.LoadArgConst("internal/"+stmt.Name.Lexeme, arg.Name, name)
	}
	wrapperSource += c.commandMapper.Call("internal/"+stmt.Name.Lexeme, "")
	wrapperSource += c.commandMapper.Return()
	c.createFunction(stmt.Name.Lexeme, wrapperSource, args, stmt.ReturnType)
	c.createFunction("internal/"+stmt.Name.Lexeme, source, args, stmt.ReturnType)
	return ""
}

func (c *Compiler) VisitBlock(stmt statements.BlockStmt) string {
	section := ""
	for _, s := range stmt.Statements {
		section += s.Accept(c)
	}
	//c.currentFunctionSections = append(c.currentFunctionSections, section)
	return section // c.commandMapper.CallSection(c.currentFunction.Name.Lexeme, len(c.currentFunctionSections)-1)
}

func (c *Compiler) VisitExpression(stmt statements.ExpressionStmt) string {
	return stmt.Expression.Accept(c)
}

func (c *Compiler) VisitReturn(stmt statements.ReturnStmt) string {
	cmd := ""

	expr := stmt.Expression
	if expr != nil {
		if expr.ReturnType() != c.currentFunction.ReturnType {
			log.Fatalln("Return type does not match function return type")
		}
		cmd += expr.Accept(c)
	}
	cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), mapping.RET)
	cmd += c.commandMapper.Return()
	return cmd
}

func (c *Compiler) VisitVariableDeclaration(stmt statements.VariableDeclarationStmt) string {
	cmd := ""
	if stmt.Initializer != nil {
		cmd += stmt.Initializer.Accept(c)
		cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(stmt.Name.Lexeme))
	} else {
		switch stmt.Type {
		case types.IntType:
			cmd += c.commandMapper.MakeConst(nbt.NewInt(0), c.commandMapper.Cs(stmt.Name.Lexeme))
		case types.DoubleType:
			cmd += c.commandMapper.MakeConst(nbt.NewDouble(0), c.commandMapper.Cs(stmt.Name.Lexeme))
		case types.StringType:
			cmd += c.commandMapper.MakeConst(nbt.NewString(""), c.commandMapper.Cs(stmt.Name.Lexeme))
		default:
			switch stmt.Type.(type) {
			case types.ListTypeStruct:
				cmd += c.commandMapper.MakeConst(nbt.NewList(), c.commandMapper.Cs(stmt.Name.Lexeme))
			case types.StructTypeStruct:
				cmd += c.commandMapper.MoveRaw(
					fmt.Sprintf("%s:data", c.Namespace),
					fmt.Sprintf("%s.%s", mapping.StructPath, stmt.Type.ToString()),
					fmt.Sprintf("%s:data", c.Namespace),
					fmt.Sprintf("%s.%s", mapping.VarPath, c.commandMapper.Cs(stmt.Name.Lexeme)),
				)
			default:
				c.error(stmt.Name.SourceLocation, fmt.Sprintf("Invalid type: %v", stmt.Type))
			}
		}
	}
	c.scopes[c.currentScope] = append(c.scopes[c.currentScope],
		TypedIdentifier{
			stmt.Name.Lexeme,
			stmt.Type,
		})
	return cmd
}

func (c *Compiler) VisitVariableAssignment(stmt statements.VariableAssignmentStmt) string {
	cmd := ""
	isIndexedAssignment := len(stmt.Accessors) > 0
	if !stmt.Value.ReturnType().Equals(c.getReturnType(stmt.Name.Lexeme)) && !isIndexedAssignment {
		c.error(stmt.Name.SourceLocation, fmt.Sprintf("Assignment type mismatch: %v != %v", c.getReturnType(stmt.Name.Lexeme).ToString(), stmt.Value.ReturnType().ToString()))
	}
	cmd += stmt.Value.Accept(c)
	valueReg := c.commandMapper.Cs(c.newRegister(mapping.RX))
	cmd += c.commandMapper.MakeConst(nbt.NewInt(0), valueReg)
	cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), valueReg)
	if isIndexedAssignment {
		pathReg := c.commandMapper.Cs(c.newRegister(mapping.RX))
		cmd += c.commandMapper.MakeConst(nbt.NewString(""), pathReg)
		for i := 0; i < len(stmt.Accessors); i++ {
			cmd += fmt.Sprintf("### BEGIN Compute path part %v/%v ###\n", i+1, len(stmt.Accessors))
			switch stmt.Accessors[i].(type) {
			case statements.IndexAccessor:
				indexAccessor := stmt.Accessors[i].(statements.IndexAccessor)
				cmd += "###       Compile index expression ###\n"
				cmd += indexAccessor.Index.Accept(c)
				cmd += "###       Move to its own register ###\n"
				indexReg := c.commandMapper.Cs(c.newRegister(mapping.RX))
				cmd += c.commandMapper.MakeConst(nbt.NewInt(0), indexReg)
				cmd += c.commandMapper.Move(c.commandMapper.Cs(mapping.RX), indexReg)
				// wrap index in brackets and append to pathReg
				// pathReg += "[" + indexReg + "]"
				cmd += "###       Wrap in brackets ###\n"
				cmd += c.commandMapper.MakeIndex(indexReg, indexReg)
				cmd += "###       Append to path ###\n"
				cmd += c.commandMapper.Concat(pathReg, indexReg, pathReg)
			case statements.FieldAccessor:
				fieldAccessor := stmt.Accessors[i].(statements.FieldAccessor)
				cmd += "###       Compile field expression ###\n"
				fieldReg := c.commandMapper.Cs(c.newRegister(mapping.RX))
				cmd += c.commandMapper.MakeConst(nbt.NewString(fieldAccessor.ToString()), fieldReg)
				cmd += "###       Append to path ###\n"
				cmd += c.commandMapper.Concat(pathReg, fieldReg, pathReg)
			}
			cmd += fmt.Sprintf("### END   Compute path part %v/%v ###\n", i+1, len(stmt.Accessors))
		}
		cmd += c.commandMapper.PathSet(c.commandMapper.Cs(stmt.Name.Lexeme), pathReg, valueReg)
	} else {
		cmd += c.commandMapper.Move(valueReg, c.commandMapper.Cs(stmt.Name.Lexeme))
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
	condVar := c.newRegister(mapping.RX)
	cmd += c.commandMapper.MakeConst(nbt.NewInt(0), condVar)
	cmd += c.commandMapper.MoveScore(c.commandMapper.Cs(mapping.RX), c.commandMapper.Cs(condVar))
	cmd += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 1", c.commandMapper.Cs(condVar), c.Namespace), true, thenSource)
	cmd += c.commandMapper.ExecCond(fmt.Sprintf("score %s %s matches 1", c.commandMapper.Cs(condVar), c.Namespace), false, elseSource)

	return cmd
}

func (c *Compiler) VisitStructDeclaration(stmt statements.StructDeclarationStmt) string {
	return c.commandMapper.StructDefine(stmt)
}
