package compiler

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/ir"
	"github.com/Kolterdyx/mcbasic/internal/paths"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	"github.com/Kolterdyx/mcbasic/internal/types"
	"github.com/Kolterdyx/mcbasic/internal/utils"
	log "github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
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
	Namespace    string
	Scope        string
	DatapackRoot string
	Config       interfaces.ProjectConfig

	registerCounter int
	storage         string

	compiledFunctions map[string]interfaces.Function
	branchCounter     int

	currentScope *symbol.Table

	libs embed.FS

	funcPath string
}

func NewCompiler(
	config interfaces.ProjectConfig,
	libs embed.FS,
) *Compiler {
	c := &Compiler{
		storage:      fmt.Sprintf("%s:data", config.Project.Namespace),
		Config:       config,
		DatapackRoot: config.OutputDir,
		Namespace:    config.Project.Namespace,
		libs:         libs,
	}

	c.Namespace = c.Config.Project.Namespace
	c.DatapackRoot, _ = filepath.Abs(path.Join(c.Config.OutputDir, c.Config.Project.Name))
	return c
}

func (c *Compiler) Compile(sourceAst ast.Source, symbols *symbol.Table) error {
	c.compiledFunctions = make(map[string]interfaces.Function)
	c.currentScope = symbols
	c.createBasePack()
	functions := c.compileToIR(sourceAst)
	functions = append(functions, c.compileBuiltins()...)
	functions = c.optimizeIRCode(functions)
	c.compileIRtoDatapack(functions)
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

func (c *Compiler) compileIRtoDatapack(ir []interfaces.Function) {
	for _, f := range ir {
		err := c.writeMcFunction(f.GetName(), f.ToMCFunction())
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func (c *Compiler) optimizeIRCode(functions []interfaces.Function) []interfaces.Function {
	optimizationPasses := 3
	for i, f := range functions {
		for j := 0; j < optimizationPasses; j++ {
			functions[i] = ir.OptimizeFunctionBody(f)
		}
	}
	return functions
}

func (c *Compiler) compileToIR(source ast.Source) []interfaces.Function {
	functions := make([]interfaces.Function, 0)
	for _, statement := range source {
		if statement.Type() == ast.FunctionDeclarationStatement {
			c.compiledFunctions = make(map[string]interfaces.Function)
			ast.AcceptStmt[interfaces.IRCode](statement, c)
			functions = append(functions, slices.Collect(maps.Values(c.compiledFunctions))...)
		}
	}
	return functions
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
	ns, fn := utils.SplitFunctionName(fullName, c.Namespace)
	return os.WriteFile(path.Join(c.getFuncPath(ns), fmt.Sprintf("%s.mcfunction", fn)), []byte(c.macroLineIdentifier(source)), 0644)
}

func (c *Compiler) writeFunctionTags() error {
	// load tag
	loadTag := interfaces.McTag{
		Values: []string{
			"gm:zzz/load",
			fmt.Sprintf("mcb:%s", path.Join(paths.Internal, "init")),
		},
	}
	tickTag := interfaces.McTag{
		Values: []string{
			fmt.Sprintf("%s:tick", c.Namespace),
		},
	}

	var tagBytes []byte
	var err error

	if err = os.MkdirAll(path.Join(c.DatapackRoot, paths.MinecraftTags, paths.Functions), 0755); err != nil {
		return err
	}
	if tagBytes, err = json.Marshal(loadTag); err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(c.DatapackRoot, paths.MinecraftTags, paths.Functions, "load.json"), tagBytes, 0644); err != nil {
		return err
	}
	if tagBytes, err = json.Marshal(tickTag); err != nil {
		return err
	}
	if err = os.WriteFile(path.Join(c.DatapackRoot, paths.MinecraftTags, paths.Functions, "tick.json"), tagBytes, 0644); err != nil {
		return err
	}
	return nil
}

func (c *Compiler) registerIRFunction(fullName string, source interfaces.IRCode, args []interfaces.TypedIdentifier, returnType types.ValueType) interfaces.Function {

	ns, fn := utils.SplitFunctionName(fullName, c.Namespace)
	f := interfaces.FunctionDefinition{
		Name:       fn,
		Args:       make([]interfaces.TypedIdentifier, 0),
		ReturnType: returnType,
	}
	for _, parameter := range args {
		f.Args = append(f.Args, interfaces.TypedIdentifier{Name: parameter.Name, Type: parameter.Type})
	}
	f.Args = append(f.Args, interfaces.TypedIdentifier{Name: "__call__", Type: types.IntType})
	return ir.NewFunction(fmt.Sprintf("%s:%s", ns, fn), source)
}

func (c *Compiler) n() interfaces.IRCode {
	return ir.NewCode(c.Namespace, c.storage)
}
