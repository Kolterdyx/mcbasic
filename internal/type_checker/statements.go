package type_checker

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	log "github.com/sirupsen/logrus"
)

func (t *TypeChecker) VisitExpression(stmt ast.ExpressionStmt) any {
	return ast.AcceptExpr[types.ValueType](stmt.Expression, t)
}

func (t *TypeChecker) VisitVariableDeclaration(stmt ast.VariableDeclarationStmt) any {
	sym, _ := t.table.Lookup(stmt.Name.Lexeme)
	if stmt.Initializer != nil {
		itype := ast.AcceptExpr[types.ValueType](stmt.Initializer, t)
		if !sym.ValueType().Equals(itype) {
			t.error(stmt, fmt.Sprintf("cannot assign %s to %s", itype.ToString(), sym.ValueType().ToString()))
			return types.VoidType
		}
	}
	return sym.ValueType()
}

func (t *TypeChecker) VisitFunctionDeclaration(stmt ast.FunctionDeclarationStmt) any {

	prevTable := t.table
	newTable, ok := t.table.GetChild(stmt.Name.Lexeme)
	if !ok {
		log.Fatalf("function %s not found", stmt.Name.Lexeme)
	}
	t.table = newTable
	defer func() {
		t.table = prevTable
	}()

	ast.AcceptStmt[any](stmt.Body, t)
	return nil
}

func (t *TypeChecker) VisitVariableAssignment(stmt ast.VariableAssignmentStmt) any {
	variableSymbol, _ := t.table.Lookup(stmt.Name.Lexeme)
	variableType := variableSymbol.ValueType()
	valueType := ast.AcceptExpr[types.ValueType](stmt.Value, t)
	if len(stmt.Accessors) > 0 {
		for _, accessor := range stmt.Accessors {
			switch accessor.(type) {
			case ast.IndexAccessor:
				if utils.IsListType(variableType) {
					variableType = variableType.(types.ListTypeStruct).ContentType
				} else {
					t.error(accessor, fmt.Sprintf("cannot index %s", variableType.ToString()))
				}
			case ast.FieldAccessor:
				if utils.IsStructType(variableType) {
					fieldAccessor := accessor.(ast.FieldAccessor)
					structType := variableType.(types.StructTypeStruct)
					fieldType, ok := structType.GetFieldType(fieldAccessor.Field.Lexeme)
					if !ok {
						t.error(accessor, fmt.Sprintf("fieldType %s not found in %s", fieldAccessor.Field.Lexeme, variableType.ToString()))
					}
					variableType = fieldType
				} else {
					t.error(accessor, fmt.Sprintf("cannot access field %s of %s", accessor.ToString(), variableType.ToString()))
				}
			}
		}
	}
	if !variableType.Equals(valueType) {
		t.error(stmt, fmt.Sprintf("type mismatch: cannot assign value of type %s to variable '%s' of type %s", valueType.ToString(), stmt.Name.Lexeme, variableType.ToString()))
	}
	return nil
}

func (t *TypeChecker) VisitStructDeclaration(stmt ast.StructDeclarationStmt) any {
	// There is nothing to check here
	return nil
}

func (t *TypeChecker) VisitBlock(stmt ast.BlockStmt) any {
	for _, s := range stmt.Statements {
		ast.AcceptStmt[any](s, t)
	}
	return nil
}

func (t *TypeChecker) VisitWhile(stmt ast.WhileStmt) any {
	// There is nothing to check here
	return nil
}

func (t *TypeChecker) VisitIf(stmt ast.IfStmt) any {
	ast.AcceptExpr[types.ValueType](stmt.Condition, t)
	ast.AcceptStmt[any](stmt.ThenBranch, t)
	if stmt.ElseBranch != nil {
		ast.AcceptStmt[any](stmt.ElseBranch, t)
	}
	return nil
}

func (t *TypeChecker) VisitReturn(stmt ast.ReturnStmt) any {
	var vtype types.ValueType = types.VoidType
	if stmt.Expression != nil {
		vtype = ast.AcceptExpr[types.ValueType](stmt.Expression, t)
	}
	sym, _ := t.table.Lookup(t.table.ScopeName())
	if !sym.ValueType().Equals(vtype) {
		t.error(stmt, fmt.Sprintf("cannot return %s from function %s", vtype.ToString(), sym.Name()))
	}
	return vtype
}

func (t *TypeChecker) VisitSetReturnFlag(_ ast.SetReturnFlagStmt) any {
	// There is nothing to check here
	return nil
}

func (t *TypeChecker) VisitImport(_ ast.ImportStmt) any {
	// There is nothing to check here
	return nil
}

func (t *TypeChecker) VisitExec(stmt ast.ExecStmt) any {
	etype := ast.AcceptExpr[types.ValueType](stmt.Expression, t)
	if etype != types.StringType {
		t.error(stmt, fmt.Sprintf("exec statement requires a string, got %s", stmt.GetResolvedType().ToString()))
	}
	return types.VoidType
}
