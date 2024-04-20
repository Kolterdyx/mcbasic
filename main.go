package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler"
	log "github.com/sirupsen/logrus"
	"os"
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

	versionPtr := flag.Bool("version", false, "Print version")

	projectFilePtr := flag.String("config", "project.toml", "Path to the project config file")
	outputDirPtr := flag.String("output", "build", "Output directory")
	enableTracesPtr := flag.Bool("traces", false, "Enable traces")
	fixedPointPrecision := flag.Int("fpp", 4, "Fixed point precision")
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
	config.FixedPointPrecision = *fixedPointPrecision

	if *fixedPointPrecision < 0 {
		log.Fatalln("Fixed point precision must be a positive integer")
	}
	if *fixedPointPrecision > 4 {
		log.Warnln("Fixed point precision greater than 4 may cause overflow")
	}

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
