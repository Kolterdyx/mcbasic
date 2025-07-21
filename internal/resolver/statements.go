package resolver

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
)

func (r *Resolver) VisitExpression(stmt ast.ExpressionStmt) any {
	return ast.AcceptExpr[Result](stmt.Expression, r)
}

func (r *Resolver) VisitVariableDeclaration(stmt ast.VariableDeclarationStmt) any {
	err := r.table.Define(symbol.NewSymbol(stmt.Name.Lexeme, symbol.VariableSymbol, stmt, stmt.ValueType))
	if err != nil {
		return r.error(stmt, fmt.Sprintf("variable %s already defined", stmt.Name.Lexeme))
	}
	if stmt.Initializer != nil {
		stmt.Initializer.Accept(r)
	}
	return nil
}

func (r *Resolver) VisitFunctionDeclaration(stmt ast.FunctionDeclarationStmt) any {

	err := r.table.Define(symbol.NewSymbol(stmt.Name.Lexeme, symbol.FunctionSymbol, stmt, stmt.ReturnType))
	if err != nil {
		return r.error(stmt, fmt.Sprintf("function %s already defined", stmt.Name.Lexeme))
	}

	newTable := symbol.NewTable(r.table, stmt.Name.Lexeme, r.table.OriginFile())
	prevTable := r.table
	r.table = newTable
	prevTable.AddChild(newTable)
	defer func() {
		r.table = prevTable
	}()

	for _, param := range stmt.Parameters {
		err := r.table.Define(symbol.NewSymbol(param.Name.Lexeme, symbol.VariableSymbol, param, param.ValueType))
		if err != nil {
			return r.error(param, fmt.Sprintf("parameter %s already defined", param.Name.Lexeme))
		}
	}
	return stmt.Body.Accept(r)
}

func (r *Resolver) VisitVariableAssignment(stmt ast.VariableAssignmentStmt) any {
	_, ok := r.table.Lookup(stmt.Name.Lexeme)
	if !ok {
		return r.error(stmt, fmt.Sprintf("variable %s not defined", stmt.Name.Lexeme))
	}
	if stmt.Value != nil {
		stmt.Value.Accept(r)
	}
	return nil
}

func (r *Resolver) VisitStructDeclaration(stmt ast.StructDeclarationStmt) any {
	err := r.table.Define(symbol.NewSymbol(stmt.Name.Lexeme, symbol.StructSymbol, stmt, stmt.StructType))
	if err != nil {
		return r.error(stmt, fmt.Sprintf("struct %s already defined", stmt.Name.Lexeme))
	}
	return nil
}

func (r *Resolver) VisitBlock(stmt ast.BlockStmt) any {
	for _, s := range stmt.Statements {
		s.Accept(r)
	}
	return nil
}

func (r *Resolver) VisitWhile(_ ast.WhileStmt) any {
	return nil
}

func (r *Resolver) VisitIf(stmt ast.IfStmt) any {
	stmt.Condition.Accept(r)
	stmt.ThenBranch.Accept(r)
	if stmt.ElseBranch != nil {
		stmt.ElseBranch.Accept(r)
	}
	return nil
}

func (r *Resolver) VisitReturn(stmt ast.ReturnStmt) any {
	if stmt.Expression != nil {
		stmt.Expression.Accept(r)
	}
	return nil
}

func (r *Resolver) VisitSetReturnFlag(_ ast.SetReturnFlagStmt) any {
	return nil
}

func (r *Resolver) VisitImport(stmt ast.ImportStmt) any {
	moduleTable, ok := r.symbolManager.GetFile(stmt.Path)
	if !ok {
		return r.error(stmt, fmt.Sprintf("module %s not found", stmt.Path))
	}
	for alias, symToken := range stmt.SymbolNames {
		symName := symToken.Lexeme
		if symName == "*" {
			err := r.table.ImportTable(moduleTable)
			if err != nil {
				return r.error(stmt, err.Error())
			}
			continue
		}
		sym, ok := moduleTable.Lookup(symName)
		if !ok {
			return r.error(stmt, fmt.Sprintf("symbol %s not found in module %s", symName, stmt.Path))
		}
		err := r.table.ImportSymbol(symbol.NewAlias(alias, sym))
		if err != nil {
			return r.error(stmt, fmt.Sprintf("symbol %s already defined", sym.Name()))
		}
	}
	return nil
}

func (r *Resolver) VisitExec(stmt ast.ExecStmt) any {
	stmt.Expression.Accept(r)
	return nil
}
