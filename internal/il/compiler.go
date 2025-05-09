package il

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
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
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
	return &Compiler{
		storage:   fmt.Sprintf("%s:data", config.Project.Namespace),
		Config:    config,
		Namespace: config.Project.Namespace,
		headers:   headers,
		libs:      libs,
	}
}

func (c *Compiler) Compile(program parser.Program) {
	c.Structs = program.Structs
	c.compiledFunctions = make(map[string]string)
	c.scopes = make(map[string][]interfaces.TypedIdentifier)

	c.createBasePack()
	c.setFunctionDefinitions(program)
	c.compileToIR(program)
	ir := c.optimizeIRCode()
	ir = c.addStructDeclarationFunction(program, ir)
	c.compileIRtoDatapack(ir)
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
	c.createPackMeta()
}

func (c *Compiler) compileIRtoDatapack(ir []Function) {
	for _, f := range ir {
		compiledFunc := f.ToMCFunction()
		c.createFunction(f.Name, compiledFunc, nil, types.VoidType)
	}
}

func (c *Compiler) addStructDeclarationFunction(program parser.Program, ir []Function) []Function {
	structDefFuncSource := ""
	for _, structDef := range program.Structs {
		structDefFuncSource += structDef.Accept(c)
	}
	structDefFuncSource = "func internal/struct_definitions\n" + structDefFuncSource + "\nret\n"
	structIr := ParseIL(structDefFuncSource, c.Namespace, c.storage)
	if len(structIr) != 1 {
		log.Fatalln("Struct definition function should have one function")
	}
	return append(ir, structIr[0])
}

func (c *Compiler) optimizeIRCode() []Function {
	ilSource := strings.Join(slices.Collect(maps.Values(c.compiledFunctions)), "\n")
	ir := ParseIL(ilSource, c.Namespace, c.storage)
	optimizationPasses := 3
	for i, f := range ir {
		for j := 0; j < optimizationPasses; j++ {
			ir[i] = OptimizeFunctionBody(f)
		}
	}
	return ir
}

func (c *Compiler) compileToIR(program parser.Program) {
	for _, f := range program.Functions {
		c.compiledFunctions[f.Name.Lexeme] = f.Accept(c)
	}
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
	c.Namespace = c.Config.Project.Namespace
	c.DatapackRoot, _ = filepath.Abs(path.Join(c.Config.OutputDir, c.Config.Project.Name))
	log.Infof("Compiling to %s\n", c.DatapackRoot)
	c.funcPath = c.getFuncPath(c.Namespace)

	errs := []error{
		os.MkdirAll(path.Join(c.funcPath, paths.Internal), 0755),
		os.MkdirAll(path.Join(c.funcPath, paths.FunctionBranches), 0755),
		os.MkdirAll(path.Join(paths.McbFunctions, paths.Internal), 0755),
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

func (c *Compiler) createFunction(fullName string, source string, args []interfaces.TypedIdentifier, returnType types.ValueType) {
	//if fullName == c.LoadFuncName || fullName == c.TickFuncName {
	//	return
	//}

	// If the function name is in the format of "namespace:function", get the namespace from the name
	if fullName == "" {
		c.error(interfaces.SourceLocation{}, "Function name cannot be empty")
		return
	}
	namespace, name := splitFunctionName(fullName, c.Namespace)
	filename := name + ".mcfunction"

	f := interfaces.FunctionDefinition{
		Name:       name,
		Args:       make([]interfaces.TypedIdentifier, 0),
		ReturnType: returnType,
	}
	for _, parameter := range args {
		f.Args = append(f.Args, interfaces.TypedIdentifier{Name: parameter.Name, Type: parameter.Type})
	}
	f.Args = append(f.Args, interfaces.TypedIdentifier{Name: "__call__", Type: types.IntType})
	c.functionDefinitions[fullName] = f

	err := os.WriteFile(c.getFuncPath(namespace)+"/"+filename, []byte(c.macroLineIdentifier(source)), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
