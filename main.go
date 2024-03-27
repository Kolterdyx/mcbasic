package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

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
	log.Debug("Source loaded successfully")

	// Create a scanner
	scanner := &internal.Scanner{}
	tokens := scanner.Scan(source)
	if scanner.HadError {
		os.Exit(1)
	}
	log.Debug("Tokens scanned successfully")
	parser_ := parser.Parser{Tokens: tokens}
	program := parser_.Parse()
	if parser_.HadError {
		os.Exit(1)
	}
	log.Debug("Program parsed successfully")
	c := compiler.Compiler{Config: config}
	c.Compile(program)
	log.Debug("Compilation complete")
}

func loadProject(file string) interfaces.ProjectConfig {
	// Read the file
	blob, _ := os.ReadFile(file)
	tomlString := string(blob)
	var project interfaces.ProjectConfig
	_, err := toml.Decode(tomlString, &project)
	if err != nil {
		log.Fatalln(err)
	}
	return project
}
