package frontend

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/resolver"
	"github.com/Kolterdyx/mcbasic/internal/scanner"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
	log "github.com/sirupsen/logrus"
	"os"
)

type Frontend struct {
	visitedFiles  map[string]bool
	units         map[string]*CompilationUnit
	symbolManager *symbol.Manager
}

func NewFrontend(rootPath string) *Frontend {
	return &Frontend{
		visitedFiles:  make(map[string]bool),
		units:         make(map[string]*CompilationUnit),
		symbolManager: symbol.NewManager(rootPath),
	}
}

func (f *Frontend) Parse(path string) error {
	if f.visitedFiles[path] {
		return nil // already parsed
	}
	f.visitedFiles[path] = true

	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	tokens, errs := scanner.Scan(string(src))
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to scan file: %s", path)
	}

	p := parser.NewParser(tokens)

	fileAst, errs := p.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to parse file: %s", path)
	}

	table := symbol.NewTable(nil, "file:"+path, path)
	r := resolver.NewResolver(fileAst, table)
	errs = r.Resolve()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to resolve file: %s", path)
	}

	t := typeChecker.NewTypeChecker()
	errs = t.Check(fileAst, table)
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
		Tokens:   tokens,
	}

	f.units[path] = unit
	f.symbolManager.AddFile(path, table)

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
