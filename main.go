package main

import (
	"flag"
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

	config := parseArgs()

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

	// Remove the contents of the output directory
	err := os.RemoveAll(config.OutputDir)
	if err != nil {
		log.Fatalln(err)
	}

	c := compiler.NewCompiler(config)
	c.Compile(program)
	log.Info("Compilation complete")
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

func parseArgs() interfaces.ProjectConfig {
	// Parse command line arguments
	projectFilePtr := flag.String("config", "project.toml", "Path to the project config file")
	outputDirPtr := flag.String("output", "build", "Output directory")
	flag.Parse()

	// Load config toml file
	config := loadProject(*projectFilePtr)
	validateProjectConfig(config)
	config.OutputDir = *outputDirPtr

	return config
}

func validateProjectConfig(config interfaces.ProjectConfig) {
	if config.Project.Entrypoint == "" {
		log.Fatalln("Entrypoint not specified")
	}
	if config.Project.Name == "" {
		log.Fatalln("Name not specified")
	}
}
