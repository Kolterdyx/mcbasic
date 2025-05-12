package frontend

import (
	"fmt"
	"github.com/Kolterdyx/mcbasic/internal/ast"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/scanner"
	"github.com/Kolterdyx/mcbasic/internal/statements"
	"github.com/Kolterdyx/mcbasic/internal/symbol"
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
			fmt.Println(err)
		}
		return fmt.Errorf("failed to scan file: %s", path)
	}

	table := symbol.NewTable(nil, "file:"+path, path)
	p := parser.NewParser(tokens, f.symbolManager, table)

	syms, fileAst, errs := p.Parse()
	if len(errs) > 0 {
		return err
	}

	unit := &CompilationUnit{
		FilePath: path,
		AST:      fileAst,
		Symbols:  syms,
		Tokens:   tokens,
	}

	f.units[path] = unit
	f.symbolManager.AddFile(path, table)

	imports := make(map[string]statements.ImportStmt)
	for _, stmt := range fileAst {
		if stmt.Type() == ast.ImportStatement {
			importStmt, ok := stmt.(statements.ImportStmt)
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
