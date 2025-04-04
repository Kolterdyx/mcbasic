package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
)

//go:embed version.txt
var version string

//go:embed libs
var libs embed.FS

//go:embed headers
var builtinHeaders embed.FS

func main() {

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

	config, projectRoot := parseArgs()

	entrypoint := config.Project.Entrypoint
	blob, _ := os.ReadFile(path.Join(projectRoot, entrypoint))
	source := string(blob)
	log.Debug("Source loaded successfully")

	// Create a scanner
	scanner := &internal.Scanner{}
	tokens := scanner.Scan(source)
	if scanner.HadError {
		os.Exit(1)
	}
	log.Debug("Tokens scanned successfully")

	headers, err := loadHeaders(config.Dependencies.Headers, projectRoot)
	if err != nil {
		log.Fatalln(err)
	}

	parser_ := parser.Parser{Tokens: tokens, Headers: headers}
	program := parser_.Parse()
	if parser_.HadError {
		os.Exit(1)
	}
	log.Debug("Program parsed successfully")

	// Remove the contents of the output directory
	err = os.RemoveAll(path.Join(config.OutputDir, config.Project.Name))
	if err != nil {
		log.Fatalln(err)
	}

	c := compiler.NewCompiler(config, projectRoot, headers, libs)
	c.Compile(program)
	log.Info("Compilation complete")
}

func loadHeaders(headerPaths []string, projectRoot string) ([]interfaces.DatapackHeader, error) {
	headers := make([]interfaces.DatapackHeader, 0)

	for i, h := range headerPaths {
		headerPath := path.Join(projectRoot, h)
		headerPaths[i] = headerPath
	}

	// include builtin headers
	builtinHeaderPaths, err := builtinHeaders.ReadDir("headers")
	if err != nil {
		return nil, err
	}
	for _, h := range builtinHeaderPaths {
		if !h.IsDir() && strings.HasSuffix(h.Name(), ".json") {
			headerPaths = append(headerPaths, path.Join("headers", h.Name()))
		}
	}

	for _, headerPath := range headerPaths {
		log.Debug("Loading header: ", headerPath)
		if _, err := os.Stat(headerPath); os.IsNotExist(err) {
			log.Warnf("Header file %s does not exist, skipping...", headerPath)
			continue
		}
		headerFile, err := os.ReadFile(headerPath)
		if err != nil {
			return nil, err
		}
		header := interfaces.DatapackHeader{}
		err = json.Unmarshal(headerFile, &header)
		if err != nil {
			return nil, err
		}
		log.Debug("Header: ", header)
		headers = append(headers, header)
	}
	log.Debugf("Headers loaded successfully")
	return headers, nil
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

func parseArgs() (interfaces.ProjectConfig, string) {
	// Parse command line arguments

	versionPtr := flag.Bool("version", false, "Print version")

	projectFilePtr := flag.String("config", "project.toml", "Path to the project config file")
	outputDirPtr := flag.String("output", "build", "Output directory")
	enableTracesPtr := flag.Bool("traces", false, "Enable traces")
	flag.Parse()

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

	return config, path.Dir(*projectFilePtr)
}

func validateProjectConfig(config interfaces.ProjectConfig) {
	if config.Project.Entrypoint == "" {
		log.Fatalln("Entrypoint not specified")
	}
	if config.Project.Name == "" {
		log.Fatalln("Name not specified")
	}
}
