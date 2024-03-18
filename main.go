package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"os"
)

func main() {
	// Read the config file
	if len(os.Args) < 2 {
		fmt.Println("Usage: mcbasic <config>")
	}
	projectFile := os.Args[1]
	if projectFile == "" {
		fmt.Println("Usage: mcbasic <config>")
		return
	}

	// Load toml file
	config := loadProject(projectFile)
	entrypoint := config.Project.Entrypoint
	blob, _ := os.ReadFile(entrypoint)
	source := string(blob)

	// Create a scanner
	scanner := &internal.Scanner{}
	tokens := scanner.Scan(source)
	if scanner.HadError {
		os.Exit(65)
	}
	parser_ := parser.Parser{Tokens: tokens}
	stmts := parser_.Parse()
	if parser_.HadError {
		os.Exit(65)
	}
	debugVisitor := internal.DebugVisitor{}
	for _, stmt := range stmts {
		stmt.Accept(debugVisitor)
	}

}

func loadProject(file string) internal.ProjectConfig {
	// Read the file
	blob, _ := os.ReadFile(file)
	tomlString := string(blob)
	var project internal.ProjectConfig
	_, err := toml.Decode(tomlString, &project)
	if err != nil {
		fmt.Println(err)
		return internal.ProjectConfig{}
	}
	return project
}
