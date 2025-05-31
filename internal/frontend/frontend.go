package frontend

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/compiler"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/resolver"
	"github.com/Kolterdyx/mcbasic/internal/scanner"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	"github.com/Kolterdyx/mcbasic/internal/tokens"
	"github.com/Kolterdyx/mcbasic/internal/type_checker"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Frontend struct {
	visitedFiles   map[string]bool
	units          map[string]*CompilationUnit
	symbolManager  *symbol.Manager
	projectRoot    string
	builtinHeaders embed.FS
	builtinLibs    embed.FS
}

func NewFrontend(projectRoot string, builtinHeaders embed.FS, builtinLibs embed.FS) *Frontend {
	f := &Frontend{
		visitedFiles:   make(map[string]bool),
		units:          make(map[string]*CompilationUnit),
		symbolManager:  symbol.NewManager(projectRoot),
		projectRoot:    projectRoot,
		builtinHeaders: builtinHeaders,
		builtinLibs:    builtinLibs,
	}
	builtinHeaderFiles, err := builtinHeaders.ReadDir("headers")
	if err != nil {
		log.Fatalf("Failed to read builtin headers: %v", err)
	}
	for _, file := range builtinHeaderFiles {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			headerPath := path.Join("headers", file.Name())
			headerFile, err := builtinHeaders.ReadFile(headerPath)
			if err != nil {
				log.Fatalf("Failed to read builtin header %s: %v", headerPath, err)
			}
			var header interfaces.DatapackHeader
			err = json.Unmarshal(headerFile, &header)
			if err != nil {
				log.Fatalf("Failed to unmarshal builtin header %s: %v", headerPath, err)
			}
			headerSymbols := symbol.NewTable(nil, header.Namespace, headerPath)
			headerFuncDefs := parser.GetHeaderFuncDefs(header)
			for funcName, funcDef := range headerFuncDefs {
				params := make([]ast.VariableDeclarationStmt, 0, len(funcDef.Parameters))
				for _, param := range funcDef.Parameters {
					params = append(params, ast.VariableDeclarationStmt{
						Name: tokens.Token{
							Lexeme:  param.Name,
							Type:    tokens.Identifier,
							Literal: param.Name,
						},
						ValueType: param.Type,
					})
				}
				err := headerSymbols.Define(
					symbol.NewSymbol(
						funcName,
						symbol.FunctionSymbol,
						ast.FunctionDeclarationStmt{
							Name: tokens.Token{
								Lexeme:  funcName,
								Type:    tokens.Identifier,
								Literal: funcName,
							},
							Parameters: params,
							ReturnType: funcDef.ReturnType,
						},
						funcDef.ReturnType,
					),
				)
				if err != nil {
					return nil
				}
			}
			f.symbolManager.AddFile(header.Namespace, headerSymbols)
			f.visitedFiles[header.Namespace] = true
			log.Debugf("Loaded builtin header: '%s'", header.Namespace)
		}
	}
	return f
}

func (f *Frontend) Parse(path string) error {
	log.Debugf("Parsing %s", path)
	if f.visitedFiles[path] {
		return nil // already parsed
	}
	f.visitedFiles[path] = true

	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	scannedTokens, errs := scanner.Scan(string(src))
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to scan file: %s", path)
	}
	relPath, err := filepath.Rel(f.projectRoot, path)
	if err != nil {
		return err
	}
	p := parser.NewParser(relPath, scannedTokens)

	fileAst, errs := p.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to parse file: %s", path)
	}

	table := symbol.NewTable(nil, "file:"+path, path)
	f.symbolManager.AddFile(path, table)
	r := resolver.NewResolver(fileAst, table, f.symbolManager)
	errs = r.Resolve()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to resolve file: %s", path)
	}

	t := type_checker.NewTypeChecker(fileAst, table)
	errs = t.Check()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to type check file: %s", path)
	}

	unit := &CompilationUnit{
		FilePath: path,
		AST:      fileAst,
		Symbols:  table,
	}

	f.units[path] = unit

	imports := make(map[string]ast.ImportStmt)
	for _, stmt := range fileAst {
		if stmt.Type() == ast.ImportStatement {
			importStmt, ok := stmt.(ast.ImportStmt)
			if !ok {
				return fmt.Errorf("failed to cast statement to ImportStatement")
			}
			imports[importStmt.Path] = importStmt
		}
	}

	// Recurse on imported files
	for _, importStmt := range imports {
		err := f.Parse(importStmt.Path)
		if err != nil {
			return err
		}

	}

	return nil
}

func (f *Frontend) Compile(config interfaces.ProjectConfig) []error {
	errs := make([]error, 0)
	for _, unit := range f.units {
		log.Debugf("Compiling %s", unit.FilePath)
		errs = append(errs, f.compileUnit(config, unit))
	}
	return errs
}

func (f *Frontend) compileUnit(config interfaces.ProjectConfig, unit *CompilationUnit) error {
	c := compiler.NewCompiler(config, f.builtinLibs)
	err := c.Compile(unit.AST, unit.Symbols)
	if err != nil {
		return err
	}
	return nil
}
