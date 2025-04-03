package compiler

import (
	"embed"
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler/ops"
	log "github.com/sirupsen/logrus"
	"io"
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
	Headers      []interfaces.DatapackHeader

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
	libs         embed.FS
}

func NewCompiler(
	config interfaces.ProjectConfig,
	projectRoot string,
	headers []interfaces.DatapackHeader,
	libs embed.FS,
) *Compiler {
	c := &Compiler{
		ProjectRoot:  projectRoot,
		Headers:      headers,
		Config:       config,
		LoadFuncName: path.Join(config.Project.Namespace, "init"),
		TickFuncName: path.Join(config.Project.Namespace, "tick"),
		libs:         libs,
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
	err = c.copyEmbeddedLibs()
	if err != nil {
		log.Fatalln(err)
	}
	c.createPackMeta()

	c.functions = utils.GetHeaderFuncDefs(c.Headers)

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

	err := os.WriteFile(c.getFuncPath(namespace)+"/"+filename, []byte(source), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *Compiler) createFunctionTags() {
	// load tag
	loadTag := `{
	"values": [
		"%s",
		"gm:zzz/load"
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

func (c *Compiler) copyEmbeddedLibs() error {
	// c.libs is a folder. It contains multiple folders that must be copied to the datapack data folder
	files, err := c.libs.ReadDir("libs")
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			// copy the folder to the datapack data folder
			err := c.copyDirRecursive(file.Name(), "libs", path.Join(c.DatapackRoot, "data"))
			if err != nil {
				return err
			}
		}
	}
	return err
}

// copyDirRecursive recursively copies a directory and its contents from source to destination
func (c *Compiler) copyDirRecursive(srcDir, baseDir, destDir string) error {
	// Read the files inside the srcDir
	files, err := c.libs.ReadDir(filepath.Join(baseDir, srcDir))
	if err != nil {
		return fmt.Errorf("error reading directory %s: %w", srcDir, err)
	}

	// Create the corresponding target directory
	destPath := filepath.Join(destDir, srcDir)
	err = os.MkdirAll(destPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory %s: %w", destPath, err)
	}

	// Iterate over the files in the source directory
	for _, file := range files {
		// Check if it's a directory or file
		if file.IsDir() {
			// Recursively copy subdirectories
			err := c.copyDirRecursive(file.Name(), filepath.Join(baseDir, srcDir), destPath)
			if err != nil {
				return err
			}
		} else {
			// Copy the file to the target location
			err := c.copyFile(file.Name(), filepath.Join(baseDir, srcDir), destPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// copyFile copies a single file from the source to the destination
func (c *Compiler) copyFile(srcFile, baseDir, destDir string) error {
	// Open the source file
	srcPath := filepath.Join(baseDir, srcFile)
	srcFileHandle, err := c.libs.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", srcPath, err)
	}
	defer srcFileHandle.Close()

	// Create the destination file
	destPath := filepath.Join(destDir, srcFile)
	destFileHandle, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating destination file %s: %w", destPath, err)
	}
	defer destFileHandle.Close()

	// Copy the file contents
	_, err = io.Copy(destFileHandle, srcFileHandle)
	if err != nil {
		return fmt.Errorf("error copying file %s: %w", srcPath, err)
	}

	return nil
}
