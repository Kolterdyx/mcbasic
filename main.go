package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/visitors"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

//go:embed version.txt
var f embed.FS

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
	err := os.RemoveAll(path.Join(config.OutputDir, config.Project.Name))
	if err != nil {
		log.Fatalln(err)
	}

	d := visitors.DebugVisitor{}
	program.Visit(d)

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

	versionPtr := flag.Bool("version", false, "Print version")

	projectFilePtr := flag.String("config", "project.toml", "Path to the project config file")
	outputDirPtr := flag.String("output", "build", "Output directory")
	enableTracesPtr := flag.Bool("traces", false, "Enable traces")
	flag.Parse()

	data, _ := f.ReadFile("version.txt")
	if data == nil {
		log.Warnln("Version file not found")
	}
	version := string(data)

	if *versionPtr {
		fmt.Printf("MCBasic version %s\n", version)
		fmt.Println("Created by Kolterdyx")
		fmt.Println("https://github.com/Kolterdyx/mcbasic")
		os.Exit(0)
	}

	// Load config toml file
	config := loadProject(*projectFilePtr)
	validateProjectConfig(config)
	config.OutputDir = *outputDirPtr
	config.EnableTraces = *enableTracesPtr

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
