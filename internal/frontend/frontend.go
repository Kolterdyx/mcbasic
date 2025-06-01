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
	"text/template"
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

func loadStdlib(f *Frontend, stdlib embed.FS) {
	stdlibFiles, err := stdlib.ReadDir("stdlib")
	if err != nil {
		log.Fatalf("Failed to read stdlib files: %v", err)
	}
	for _, filepath := range stdlibFiles {
		if !filepath.IsDir() && strings.HasSuffix(filepath.Name(), ".mcb") {
			log.Debugf("Adding stdlib %s", filepath.Name())
			file, err := stdlib.ReadFile(path.Join("stdlib", filepath.Name()))
			if err != nil {
				log.Fatalf("Failed to read stdlib file %s: %v", filepath.Name(), err)
			}
			tmpl, err := template.New("stdlib").Parse(string(file))
			if err != nil {
				log.Fatalf("Failed to parse stdlib file %s: %v", filepath.Name(), err)
			}
			var content strings.Builder
			err = tmpl.Execute(&content, f.config.Project)
			if err != nil {
				log.Fatalf("Failed to execute template for stdlib file %s: %v", filepath.Name(), err)
			}
			f.stdlib[filepath.Name()] = content.String()
		}
	}
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
	scannedTokens, errs := scanner.Scan(path, string(src))
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to scan file: %s", path)
	}
	p := parser.NewParser(path, scannedTokens)

	fileAst, errs := p.Parse()
	if len(errs) > 0 {
		for _, err := range errs {
			log.Error(err)
		}
		return fmt.Errorf("failed to parse file: %s", path)
	}

	table := symbol.NewTable(nil, "file:"+path, path)
	f.symbolManager.AddFile(path, table)

	err = f.recurseOnImportedFiles(fileAst)
	if err != nil {
		return err
	}

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
