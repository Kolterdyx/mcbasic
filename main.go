package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler"
	"os"
)

func main() {
	// Read the config file
	if len(os.Args) < 2 {
		fmt.Println("Usage: mcbasic <config>")
		return
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
	program := parser_.Parse()
	if parser_.HadError {
		os.Exit(65)
	}
	c := compiler.Compiler{Config: config}
	c.Compile(program)
}

func loadProject(file string) interfaces.ProjectConfig {
	// Read the file
	blob, _ := os.ReadFile(file)
	tomlString := string(blob)
	var project interfaces.ProjectConfig
	_, err := toml.Decode(tomlString, &project)
	if err != nil {
		fmt.Println(err)
		return interfaces.ProjectConfig{}
	}
	return project
}
