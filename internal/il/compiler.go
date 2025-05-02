package il

import (
	"embed"
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"os"
	"strings"
)

/*

IlCompiler is used to convert the AST into a simplified version of the target command, so that
it can be optimized and then converted to the final command.

*/

const (
	RX = "$RX"
	RA = "$RA"
	RB = "$RB"

	RET  = "$RET"  // Function return value
	RETF = "$RETF" // Early return flag
	CALL = "$CALL"

	VarPath    = "vars"
	ArgPath    = "args"
	StructPath = "structs"
)

type Compiler struct {
	expressions.ExprVisitor
	statements.StmtVisitor

	Namespace string
	Scope     string
	Structs   map[string]statements.StructDeclarationStmt

	RX   string
	RA   string
	RB   string
	RET  string
	CALL string

	VarPath    string
	ArgPath    string
	StructPath string

	registerCounter int
	storage         string

	compiledFunctions map[string]string
	branchCounter     int

	functions    map[string]interfaces.FuncDef
	scopes       map[string][]interfaces.TypedIdentifier
	currentScope string

	libs    embed.FS
	headers []interfaces.DatapackHeader
}

func NewCompiler(
	namespace string,
	structs map[string]statements.StructDeclarationStmt,
	headers []interfaces.DatapackHeader,
	libs embed.FS,
) *Compiler {
	return &Compiler{
		storage:   fmt.Sprintf("%s:data", namespace),
		Namespace: namespace,
		Structs:   structs,
		headers:   headers,
		libs:      libs,
	}
}

func (c *Compiler) Compile(program parser.Program) {
	c.compiledFunctions = make(map[string]string)
	c.functions = parser.GetHeaderFuncDefs(c.headers)
	c.scopes = make(map[string][]interfaces.TypedIdentifier)

	for _, function := range program.Functions {
		f := interfaces.FuncDef{
			Name:       function.Name.Lexeme,
			Args:       make([]interfaces.TypedIdentifier, 0),
			ReturnType: function.ReturnType,
		}
		for _, parameter := range function.Parameters {
			f.Args = append(f.Args, interfaces.TypedIdentifier{
				Name: parameter.Name,
				Type: parameter.Type,
			})
		}
		f.Args = append(f.Args, interfaces.TypedIdentifier{
			Name: "__call__",
			Type: types.IntType,
		})
		c.functions[function.Name.Lexeme] = f
	}

	for _, f := range program.Functions {
		c.compiledFunctions[f.Name.Lexeme] = f.Accept(c)
	}
	file, err := os.OpenFile("compiled.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	if err != nil {
		panic(err)
	}
	ilSource := ""
	for _, f := range c.compiledFunctions {
		ilSource += f
	}
	_, err = file.WriteString(strings.TrimSpace(ilSource))
	if err != nil {
		panic(err)
	}
	ir := ParseIL(ilSource)
	optimizationPasses := 3
	for i, f := range ir {
		for j := 0; j < optimizationPasses; j++ {
			ir[i] = OptimizeFunctionBody(f)
		}
	}
	optimizedIL := ""
	for _, f := range ir {
		optimizedIL += f.ToString()
	}
	optimizedFile, err := os.OpenFile("optimized.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(optimizedFile)
	if err != nil {
		panic(err)
	}
	_, err = optimizedFile.WriteString(strings.TrimSpace(optimizedIL))
	if err != nil {
		panic(err)
	}
}

func (c *Compiler) splitFunctionName(lexeme string) (string, string) {
	parts := strings.Split(lexeme, ":")
	if len(parts) == 1 {
		return c.Namespace, parts[0]
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
