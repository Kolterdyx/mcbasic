package frontend

import (
	"embed"
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/compiler"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/resolver"
	"github.com/Kolterdyx/mcbasic/internal/scanner"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	"github.com/Kolterdyx/mcbasic/internal/type_checker"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

type Frontend struct {
	visitedFiles  map[string]bool
	units         map[string]*CompilationUnit
	symbolManager *symbol.Manager
	projectRoot   string
	embedded      embed.FS
	config        interfaces.ProjectConfig
	stdlib        map[string]string // map[path]content
}

func NewFrontend(config interfaces.ProjectConfig, projectRoot string, stdlib embed.FS, embedded embed.FS) *Frontend {
	f := &Frontend{
		visitedFiles:  make(map[string]bool),
		units:         make(map[string]*CompilationUnit),
		symbolManager: symbol.NewManager(projectRoot),
		projectRoot:   projectRoot,
		embedded:      embedded,
		config:        config,
		stdlib:        make(map[string]string),
	}
	loadStdlib(f, stdlib)
	return f
}

func (f *Frontend) Parse(filePath string) error {
	log.Debugf("Parsing %s", filePath)
	if f.visitedFiles[filePath] {
		return nil // already parsed
	}
	f.visitedFiles[filePath] = true

	var src string
	if strings.HasPrefix(filePath, "@") {
		// If the filePath starts with '@', it is a stdlib file
		src2, ok := f.stdlib[filePath]
		if !ok {
			return fmt.Errorf("stdlib file %s not found", filePath)
		}
		src = src2
	} else {
		filePath, _ := strings.CutSuffix(path.Join(f.projectRoot, filePath), ".mcb")
		src2, err := os.ReadFile(filePath + ".mcb")
		if err != nil {
			return err
		}
		src = string(src2)
	}
	scannedTokens, errs := scanner.Scan(filePath, src)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to scan file: %s", filePath)
	}
	p := parser.NewParser(filePath, scannedTokens)

	fileAst, errs := p.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to parse file: %s", filePath)
	}

	table := symbol.NewTable(nil, filePath, filePath)
	f.symbolManager.AddFile(filePath, table)

	err := f.recurseOnImportedFiles(fileAst)
	if err != nil {
		return err
	}

	r := resolver.NewResolver(fileAst, table, f.symbolManager)
	errs = r.Resolve()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to resolve file: %s", filePath)
	}

	t := type_checker.NewTypeChecker(fileAst, table)
	errs = t.Check()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to type check file: %s", filePath)
	}

	unit := &CompilationUnit{
		FilePath: filePath,
		AST:      fileAst,
		Symbols:  table,
	}

	f.units[filePath] = unit
	return nil
}

func (f *Frontend) recurseOnImportedFiles(fileAst ast.Source) error {
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

	for _, importStmt := range imports {
		err := f.Parse(importStmt.Path)
		if err != nil {
			return err
		}

	}
	return nil
}

func (f *Frontend) Compile() []error {
	errs := make([]error, 0)
	for _, unit := range f.units {
		log.Debugf("Compiling %s", unit.FilePath)
		errs = append(errs, f.compileUnit(f.config, unit))
	}
	return errs
}

func (f *Frontend) compileUnit(config interfaces.ProjectConfig, unit *CompilationUnit) error {
	c := compiler.NewCompiler(config, f.embedded)
	err := c.Compile(unit.AST, unit.Symbols)
	if err != nil {
		return err
	}
	return nil
}
