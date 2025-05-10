package ir

import (
	"embed"
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/expressions"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/types"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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

	MaxCallCounter = 65536
)

type Compiler struct {
	expressions.ExprVisitor
	statements.StmtVisitor

	Namespace    string
	Scope        string
	DatapackRoot string
	Config       interfaces.ProjectConfig
	Structs      map[string]statements.StructDeclarationStmt

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

	functionDefinitions map[string]interfaces.FunctionDefinition
	scopes              map[string][]interfaces.TypedIdentifier
	currentScope        string

	libs    embed.FS
	headers []interfaces.DatapackHeader

	funcPath string
}

func NewCompiler(
	config interfaces.ProjectConfig,
	headers []interfaces.DatapackHeader,
	libs embed.FS,
) *Compiler {
	c := &Compiler{
		storage:      fmt.Sprintf("%s:data", config.Project.Namespace),
		Config:       config,
		DatapackRoot: config.OutputDir,
		Namespace:    config.Project.Namespace,
		headers:      headers,
		libs:         libs,
	}

	c.Namespace = c.Config.Project.Namespace
	c.DatapackRoot, _ = filepath.Abs(path.Join(c.Config.OutputDir, c.Config.Project.Name))
	return c
}

func (c *Compiler) Compile(program parser.Program) error {
	c.Structs = program.Structs
	c.compiledFunctions = make(map[string]string)
	c.scopes = make(map[string][]interfaces.TypedIdentifier)

	c.createBasePack()
	c.setFunctionDefinitions(program)
	ir := c.compileToIR(program)
	ir = append(ir, c.compileBuiltins()...)
	ir = c.optimizeIRCode(ir)
	ir = c.addStructDeclarationFunction(program, ir)
	c.compileIRtoDatapack(ir)
	err := c.writeFunctionTags()
	if err != nil {
		return err
	}
	return nil
}

func (c *Compiler) createBasePack() {
	err := c.createDirectoryTree()
	if err != nil {
		log.Fatalln(err)
	}
	err = c.copyEmbeddedLibs()
	if err != nil {
		log.Fatalln(err)
	}
	c.writePackMcMeta()
}

func (c *Compiler) compileIRtoDatapack(ir []Function) {
	for _, f := range ir {
		err := c.writeMcFunction(f.Name, f.ToMCFunction())
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (c *Compiler) addStructDeclarationFunction(program parser.Program, ir []Function) []Function {
	structDefFuncSource := ""
	for _, structDef := range program.Structs {
		structDefFuncSource += structDef.Accept(c)
	}
	structDefFuncSource = c.Func("internal/struct_definitions") + structDefFuncSource + c.Ret()
	structIr := ParseIRCode(structDefFuncSource, c.Namespace, c.storage)
	if len(structIr) != 1 {
		log.Fatalln("Struct definition function should have one function")
	}
	return append(ir, structIr[0])
}

func (c *Compiler) optimizeIRCode(ir []Function) []Function {
	optimizationPasses := 3
	for i, f := range ir {
		for j := 0; j < optimizationPasses; j++ {
			ir[i] = OptimizeFunctionBody(f)
		}
	}
	return ir
}

func (c *Compiler) compileToIR(program parser.Program) []Function {
	ir := make([]Function, 0)
	for _, f := range program.Functions {
		c.compiledFunctions = make(map[string]string)
		c.compiledFunctions[f.Name.Lexeme] = f.Accept(c)
		for _, funcSource := range c.compiledFunctions {
			ir = append(ir, ParseIRCode(funcSource, c.Namespace, c.storage)...)
		}
	}
	return ir
}

func (c *Compiler) setFunctionDefinitions(program parser.Program) {
	c.functionDefinitions = parser.GetHeaderFuncDefs(c.headers)
	for _, function := range program.Functions {
		f := interfaces.FunctionDefinition{
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
		c.functionDefinitions[function.Name.Lexeme] = f
	}
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

func (c *Compiler) copyDirRecursive(srcDir, baseDir, destDir string) error {
	// Read the files inside the srcDir
	files, err := c.libs.ReadDir(filepath.Join(baseDir, srcDir))
	if err != nil {
		return fmt.Errorf("error reading directory %s: %w", srcDir, err)
	}

	// Create the corresponding target directory
	destPath := path.Join(destDir, srcDir)
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

func (c *Compiler) copyFile(srcFile, baseDir, destDir string) error {
	// Open the source file
	srcPath := filepath.Join(baseDir, srcFile)
	srcFileHandle, err := c.libs.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %w", srcPath, err)
	}
	defer func(srcFileHandle fs.File) {
		err := srcFileHandle.Close()
		if err != nil {
			log.Errorf("error closing source file %s: %v", srcPath, err)
		}
	}(srcFileHandle)

	// Create the destination file
	destPath := filepath.Join(destDir, srcFile)
	destFileHandle, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating destination file %s: %w", destPath, err)
	}
	defer func(destFileHandle *os.File) {
		err := destFileHandle.Close()
		if err != nil {
			log.Errorf("error closing destination file %s: %v", destPath, err)
		}
	}(destFileHandle)

	// Copy the file contents
	_, err = io.Copy(destFileHandle, srcFileHandle)
	if err != nil {
		return fmt.Errorf("error copying file %s: %w", srcPath, err)
	}

	return nil
}

func (c *Compiler) createDirectoryTree() error {
	log.Infof("Compiling to %s\n", c.DatapackRoot)
	c.funcPath = c.getFuncPath(c.Namespace)

	errs := []error{
		os.MkdirAll(path.Join(c.funcPath, paths.Internal), 0755),
		os.MkdirAll(path.Join(c.funcPath, paths.FunctionBranches), 0755),
		os.MkdirAll(path.Join(c.DatapackRoot, paths.McbFunctions, paths.Internal), 0755),
		os.MkdirAll(path.Join(c.DatapackRoot, paths.MinecraftTags), 0755),
		os.MkdirAll(path.Join(c.DatapackRoot, paths.MathFunctions), 0755),
	}
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Compiler) writePackMcMeta() {
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

func (c *Compiler) writeMcFunction(fullName, source string) error {
	ns, fn := splitFunctionName(fullName, c.Namespace)
	return os.WriteFile(path.Join(c.getFuncPath(ns), fmt.Sprintf("%s.mcfunction", fn)), []byte(c.macroLineIdentifier(source)), 0644)
}

func (c *Compiler) writeFunctionTags() error {
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
	err := os.MkdirAll(path.Join(c.DatapackRoot, paths.MinecraftTags, paths.Functions), 0755)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(c.DatapackRoot, paths.MinecraftTags, paths.Functions, "load.json"), []byte(fmt.Sprintf(loadTag, "mcb:internal/init")), 0644)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(c.DatapackRoot, paths.MinecraftTags, paths.Functions, "tick.json"), []byte(fmt.Sprintf(tickTag, fmt.Sprintf("%s:tick", c.Namespace))), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Compiler) createIRFunction(fullName, source string, args []interfaces.TypedIdentifier, returnType types.ValueType) Function {
	if fullName == "" {
		c.error(interfaces.SourceLocation{}, "Function fn cannot be empty")
		return Function{}
	}
	ns, fn := splitFunctionName(fullName, c.Namespace)

	f := interfaces.FunctionDefinition{
		Name:       fn,
		Args:       make([]interfaces.TypedIdentifier, 0),
		ReturnType: returnType,
	}
	for _, parameter := range args {
		f.Args = append(f.Args, interfaces.TypedIdentifier{Name: parameter.Name, Type: parameter.Type})
	}
	f.Args = append(f.Args, interfaces.TypedIdentifier{Name: "__call__", Type: types.IntType})
	c.functionDefinitions[fullName] = f
	return ParseIRCode(c.Func(fmt.Sprintf("%s:%s", ns, fn))+source, c.Namespace, c.storage)[0]
}
