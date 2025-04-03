package compiler

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type TypedIdentifier struct {
	Name string
	Type interfaces.ValueType
}

type Compiler struct {
	Config       interfaces.ProjectConfig
	ProjectRoot  string
	Namespace    string
	DatapackRoot string

	mcbFuncPath string
	funcPath    string
	tagsPath    string

	currentFunction statements.FunctionDeclarationStmt
	currentScope    string

	functions map[string]interfaces.FuncDef

	scope map[string][]TypedIdentifier

	opHandler ops.Op

	expressions.ExprVisitor
	statements.StmtVisitor

	regCounters map[string]int

	LoadFuncName string
	TickFuncName string
}

func NewCompiler(config interfaces.ProjectConfig, projectRoot string) *Compiler {
	c := &Compiler{
		ProjectRoot:  projectRoot,
		Config:       config,
		LoadFuncName: path.Join(config.Project.Namespace, "init"),
		TickFuncName: path.Join(config.Project.Namespace, "tick"),
	}
	c.Namespace = config.Project.Namespace
	c.opHandler = ops.Op{
		Namespace:    c.Namespace,
		EnableTraces: config.EnableTraces,
	}
	c.functions = make(map[string]interfaces.FuncDef)
	c.scope = make(map[string][]TypedIdentifier)
	c.regCounters = make(map[string]int)

	return c
}

func (c *Compiler) Compile(program parser.Program) {
	err := c.createDirectoryTree()
	if err != nil {
		log.Fatalln(err)
	}
	c.createPackMeta()

	c.declareFunctionsFromHeaders()

	for _, function := range program.Functions {
		f := interfaces.FuncDef{
			Name:       function.Name.Lexeme,
			Args:       make([]interfaces.FuncArg, 0),
			ReturnType: function.ReturnType,
		}
		for _, parameter := range function.Parameters {
			f.Args = append(f.Args, interfaces.FuncArg{
				Name: parameter.Name,
				Type: parameter.Type,
			})
		}
		f.Args = append(f.Args, interfaces.FuncArg{
			Name: "__call__",
			Type: expressions.IntType,
		})
		c.functions[function.Name.Lexeme] = f
	}

	// Built-in functions are protected by the compiler, so they can't be overwritten
	c.createFunctionTags()
	c.createBuiltinFunctions()

	// Traverse the AST to generate the functions
	for _, f := range program.Functions {
		f.Accept(c)
	}
}

func (c *Compiler) createDirectoryTree() error {
	c.Namespace = c.Config.Project.Namespace
	c.DatapackRoot, _ = filepath.Abs(path.Join(c.Config.OutputDir, c.Config.Project.Name))
	log.Infof("Compiling to %s\n", c.DatapackRoot)
	c.funcPath = c.getFuncPath(c.Namespace)
	c.mcbFuncPath = c.getFuncPath("mcb")
	c.tagsPath = c.DatapackRoot + "/data/minecraft/tags"

	errs := []error{
		os.MkdirAll(c.funcPath, 0755),
		os.MkdirAll(c.funcPath+"/internal", 0755),
		os.MkdirAll(c.mcbFuncPath, 0755),
		os.MkdirAll(c.mcbFuncPath+"/internal", 0755),
		os.MkdirAll(c.tagsPath, 0755),
	}
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) createPackMeta() {
	packMcmeta := fmt.Sprintf(`{
	"pack": {
		"description": "%s",
		"pack_format": 71
	},
	"meta": {
		"name": "%s",
		"version": "%s"
	}
}`, c.Config.Project.Description, c.Config.Project.Name, c.Config.Project.Version)
	err := os.WriteFile(c.DatapackRoot+"/pack.mcmeta", []byte(packMcmeta), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Compiler) createBuiltinFunctions() {
	c.createFunction(
		"mcb:print",
		`$tellraw @a {text:'$(text)'}`,
		[]interfaces.FuncArg{
			{Name: "text", Type: expressions.StringType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:log",
		`$tellraw @a[tag=mcblog] {text:'$(text)',color:dark_gray,italic:true}`,
		[]interfaces.FuncArg{
			{Name: "text", Type: expressions.StringType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:exec",
		`$execute run $(command)`,
		[]interfaces.FuncArg{
			{Name: "command", Type: expressions.StringType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:internal/concat",
		`$data modify storage $(storage) $(res) set value "$(a)$(b)"`,
		[]interfaces.FuncArg{
			{Name: "storage", Type: expressions.StringType},
			{Name: "res", Type: expressions.StringType},
			{Name: "a", Type: expressions.StringType},
			{Name: "b", Type: expressions.VoidType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:internal/slice",
		`$data modify storage $(storage) $(res) set string storage $(storage) $(from) $(start) $(end)`,
		[]interfaces.FuncArg{
			{Name: "storage", Type: expressions.StringType},
			{Name: "res", Type: expressions.StringType},
			{Name: "from", Type: expressions.StringType},
			{Name: "start", Type: expressions.IntType},
			{Name: "end", Type: expressions.IntType},
		},
		expressions.VoidType,
	)
	c.createFunction(
		"mcb:len",
		fmt.Sprintf("$data modify storage %s:%s %s set value \"$(from)\"\n", c.Namespace, ops.VarPath, ops.RET)+
			fmt.Sprintf("execute store result storage %s:%s %s int 1 run data get storage %s:%s %s\n", c.Namespace, ops.VarPath, ops.RET, c.Namespace, ops.VarPath, ops.RET),
		[]interfaces.FuncArg{
			{Name: "from", Type: expressions.StringType},
		},
		expressions.IntType,
	)
	c.createFunction(
		"mcb:internal/init",
		fmt.Sprintf("scoreboard objectives add %s dummy\n", c.Namespace)+
			c.opHandler.MoveConst("0", ops.CALL)+
			c.opHandler.MoveScore(ops.CALL, ops.CALL)+
			c.opHandler.LoadArgConst("print", "text", "MCB pack loaded")+
			c.opHandler.Call("print", "")+
			c.opHandler.Call("main", ""),
		[]interfaces.FuncArg{},
		expressions.VoidType,
	)
	c.createFunction(
		"internal/tick",
		c.opHandler.Call("tick", ""),
		[]interfaces.FuncArg{},
		expressions.VoidType,
	)
}

func (c *Compiler) createFunction(fullName string, source string, args []interfaces.FuncArg, returnType interfaces.ValueType) {
	if fullName == c.LoadFuncName || fullName == c.TickFuncName {
		return
	}

	// If the function name is in the format of "namespace:function", get the namespace from the name
	if fullName == "" {
		c.error(interfaces.SourceLocation{}, "Function name cannot be empty")
		return
	}
	namespace := c.Namespace
	name := fullName
	if strings.Contains(fullName, ":") {
		parts := strings.Split(fullName, ":")
		if len(parts) != 2 {
			c.error(interfaces.SourceLocation{}, "Invalid function name format")
			return
		}
		name = parts[1]
		namespace = parts[0]
	}
	filename := name + ".mcfunction"

	f := interfaces.FuncDef{
		Name:       name,
		Args:       make([]interfaces.FuncArg, 0),
		ReturnType: returnType,
	}
	for _, parameter := range args {
		f.Args = append(f.Args, interfaces.FuncArg{Name: parameter.Name, Type: parameter.Type})
	}
	f.Args = append(f.Args, interfaces.FuncArg{Name: "__call__", Type: expressions.IntType})
	c.functions[fullName] = f
	log.Debugf("Creating function: %s", fullName)

	err := os.WriteFile(c.getFuncPath(namespace)+"/"+filename, []byte(source), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Compiler) createFunctionTags() {
	// load tag
	loadTag := `{
	"values": [
		"%s"
	]
}`
	tickTag := `{
	"values": [
		"%s"
	]
}`
	err := os.MkdirAll(c.tagsPath+"/function", 0755)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(c.tagsPath+"/function/load.json", []byte(fmt.Sprintf(loadTag, c.Namespace+":internal/init")), 0644)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.WriteFile(c.tagsPath+"/function/tick.json", []byte(fmt.Sprintf(tickTag, c.Namespace+":internal/tick")), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Compiler) error(location interfaces.SourceLocation, message string) {
	log.Errorf("[Position %d:%d] Error at '%s':\n", location.Row+1, location.Col+1, message)
}

func (c *Compiler) newRegister(regName string) string {
	c.regCounters[regName]++
	return regName + fmt.Sprintf("_%d", c.regCounters[regName])
}

// Searches the current scope for functions and variables, returns the type of the variable or function
func (c *Compiler) getReturnType(name string) interfaces.ValueType {
	for _, identifier := range c.scope[c.currentScope] {
		if identifier.Name == name {
			return identifier.Type
		}
	}
	return expressions.VoidType
}

func (c *Compiler) Compare(expr expressions.BinaryExpr, ra string, rb string, rx string) string {
	cmd := ""
	cmd += "### Comparison operation ###\n"
	switch expr.Operator.Type {
	case tokens.EqualEqual:
		if expr.Left.ReturnType() != expr.Right.ReturnType() {
			// Return false
			cmd += c.opHandler.MoveConst("0", rx)
		} else {
			if expr.Left.ReturnType() == expressions.IntType {
				cmd += c.opHandler.EqNumbers(ra, rb, rx)
			} else if expr.Left.ReturnType() == expressions.StringType {
				cmd += c.opHandler.EqStrings(ra, rb, rx)
			}
		}
	case tokens.BangEqual:
		if expr.Left.ReturnType() != expr.Right.ReturnType() {
			// Return true
			cmd += c.opHandler.MoveConst("1", rx)
		} else {
			if expr.Left.ReturnType() == expressions.IntType {
				cmd += c.opHandler.NeqNumbers(ra, rb, rx)
			} else if expr.Left.ReturnType() == expressions.StringType {
				cmd += c.opHandler.NeqStrings(ra, rb, rx)
			}

		}
	case tokens.Greater:
		if expr.Left.ReturnType() != expressions.IntType {
			c.error(expr.SourceLocation, "Invalid type in binary operation")
		}
		cmd += c.opHandler.GtNumbers(ra, rb, rx)
	case tokens.GreaterEqual:
		if expr.Left.ReturnType() != expressions.IntType {
			c.error(expr.SourceLocation, "Invalid type in binary operation")
		}
		cmd += c.opHandler.GteNumbers(ra, rb, rx)
	case tokens.Less:
		if expr.Left.ReturnType() != expressions.IntType {
			c.error(expr.SourceLocation, "Invalid type in binary operation")
		}
		cmd += c.opHandler.LtNumbers(ra, rb, rx)
	case tokens.LessEqual:
		if expr.Left.ReturnType() != expressions.IntType {
			c.error(expr.SourceLocation, "Invalid type in binary operation")
		}
		cmd += c.opHandler.LteNumbers(ra, rb, rx)
	default:
		c.error(expr.SourceLocation, "Unknown comparison operator")
	}
	return cmd
}

func (c *Compiler) getFuncPath(namespace string) string {
	return fmt.Sprintf("%s/data/%s/function", c.DatapackRoot, namespace)
}

func (c *Compiler) declareFunctionsFromHeaders() {
	for _, header := range c.Config.Dependencies.Headers {
		headerPath := path.Join(c.ProjectRoot, header)
		if _, err := os.Stat(headerPath); os.IsNotExist(err) {
			log.Warnf("Header file %s does not exist, skipping...", header)
			continue
		}
		headerFile, err := os.ReadFile(headerPath)
		if err != nil {
			log.Fatalln(err)
		}
		lines := strings.Split(string(headerFile), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "function") {
				name := strings.TrimPrefix(line, "function")
				name = strings.TrimSpace(name)
				c.functions[name] = interfaces.FuncDef{Name: name}
			}
		}
	}
}
