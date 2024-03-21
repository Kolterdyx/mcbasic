package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/visitors"
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
	block := parser_.Parse()
	if parser_.HadError {
		os.Exit(65)
	}
	visitor := visitors.DebugVisitor{}
	fmt.Println(block.Accept(visitor))

	visitorJson := visitors.JsonVisitor{}
	res := block.Accept(visitorJson)
	b, _ := json.MarshalIndent(res, "", "  ")
	// Write the json to a file
	_ = os.WriteFile("output.json", b, 0644)
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
