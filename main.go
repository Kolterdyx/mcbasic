package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"os"
)

func main() {

	// Read the config file
	projectFile := os.Args[1]
	if projectFile == "" {
		fmt.Println("Usage: mcbasic <config-file>")
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
	for _, token := range tokens {
		var literal string = ""
		if token.Literal != nil {
			literal = fmt.Sprintf("%v", token.Literal)
		}
		fmt.Printf("%-15s %-20s %d\n", token.Type, literal, token.Line)
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
