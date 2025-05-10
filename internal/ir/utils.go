package ir

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
	"path"
	"strings"
)

func (c *Compiler) inst(instruction InstructionType, args ...string) string {
	return fmt.Sprintf("%s %s\n", instruction, strings.Join(args, " "))
}

func (c *Compiler) error(location interfaces.SourceLocation, message string) {
	log.Errorf("[Position %s] Exception: %s\n", location.ToString(), message)
}

func (c *Compiler) varPath(path string) string {
	if strings.HasPrefix(path, fmt.Sprintf("%s.", VarPath)) {
		return path
	}
	return fmt.Sprintf("%s.%s", VarPath, path)
}

func (c *Compiler) argPath(funcName, arg string) string {
	return fmt.Sprintf("%s.%s.%s", ArgPath, funcName, arg)
}

func (c *Compiler) structPath(path string) string {
	return fmt.Sprintf("%s.%s", StructPath, path)
}

func (c *Compiler) makeReg(reg string) string {
	c.registerCounter++
	return fmt.Sprintf("%s%d", reg, c.registerCounter)
}

// Searches the current scopes for functionDefinitions and variables, returns the type of the variable or function
func (c *Compiler) getReturnType(name string) types.ValueType {
	for _, identifier := range c.scopes[c.currentScope] {
		if identifier.Name == name {
			return identifier.Type
		}
	}
	return types.VoidType
}

func (c *Compiler) macroLineIdentifier(source string) string {
	lines := strings.Split(source, "\n")
	if len(lines) == 0 {
		return ""
	}
	for i, line := range lines {
		if strings.Contains(line, "$(") && !(line[0:1] == "$") {
			lines[i] = "$" + line
		}
	}
	return strings.Join(lines, "\n")
}

func (c *Compiler) cmpOperator(operator tokens.TokenType) string {
	switch operator {
	case tokens.Greater:
		return ">"
	case tokens.GreaterEqual:
		return ">="
	case tokens.Less:
		return "<"
	case tokens.LessEqual:
		return "<="
	case tokens.EqualEqual:
		return "="
	case tokens.BangEqual:
		return "!="
	default:
	}
	log.Fatalln("unknown operator")
	return ""
}

func (c *Compiler) getFuncPath(namespace string) string {
	return path.Join(c.DatapackRoot, paths.Data, namespace, paths.Functions)
}

func splitFunctionName(lexeme, namespace string) (string, string) {
	parts := strings.Split(lexeme, ":")
	if len(parts) == 1 {
		return namespace, parts[0]
	}
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	panic(fmt.Sprintf("Invalid function name: %s", lexeme))
}

func (c *Compiler) makeBranchFunction(branchName string, body statements.BlockStmt) statements.FunctionDeclarationStmt {

	// If the body contains a return statement, add a statement before to set the RETF flag
	for i, stmt := range body.Statements {
		if _, ok := stmt.(statements.ReturnStmt); ok {
			body.Statements = append(body.Statements[:i], append([]statements.Stmt{
				statements.VariableDeclarationStmt{
					Name: tokens.Token{
						Type:   tokens.Identifier,
						Lexeme: RETF,
					},
					Type: types.IntType,
					Initializer: expressions.LiteralExpr{
						Value: "1",
						SourceLocation: interfaces.SourceLocation{
							Row: 0,
							Col: 0,
						},
						ValueType: types.IntType,
					},
				},
			}, body.Statements[i:]...)...)
			break
		}
	}
	return statements.FunctionDeclarationStmt{
		Name: tokens.Token{
			Type:   tokens.Identifier,
			Lexeme: branchName,
		},
		Body: body,
	}
}
