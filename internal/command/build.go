package command

import (
	"embed"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	"github.com/Kolterdyx/mcbasic/internal/parser"
	"github.com/Kolterdyx/mcbasic/internal/visitors/compiler"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"os"
	"path"
	"strings"
)

func Build(cmd *cli.Command, builtinHeaders, libs embed.FS) error {
	if cmd.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	config, projectRoot := parseArgs(cmd)

	entrypoint := config.Project.Entrypoint
	blob, _ := os.ReadFile(path.Join(projectRoot, entrypoint))
	source := string(blob)
	log.Debug("Source loaded successfully")

	// Create a scanner
	scanner := &internal.Scanner{}
	tokens := scanner.Scan(source)
	if scanner.HadError {
		return cli.Exit("There was an error parsing the program", 1)
	}
	log.Debug("Tokens scanned successfully")

	headers, err := loadHeaders(config.Dependencies.Headers, projectRoot, builtinHeaders)
	if err != nil {
		return err
	}

	parser_ := parser.Parser{Tokens: tokens, Headers: headers}
	program := parser_.Parse()
	if parser_.HadError {
		return cli.Exit("There was an error parsing the program", 1)
	}
	log.Debug("Program parsed successfully")

	// Remove the contents of the output directory
	err = os.RemoveAll(path.Join(config.OutputDir, config.Project.Name))
	if err != nil {
		return err
	}

	c := compiler.NewCompiler(config, projectRoot, headers, libs)
	c.Compile(program)
	log.Info("Compilation complete")
	return nil
}

func loadHeaders(headerPaths []string, projectRoot string, builtinHeaders embed.FS) ([]interfaces.DatapackHeader, error) {
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

func parseArgs(cmd *cli.Command) (interfaces.ProjectConfig, string) {

	projectFile := cmd.String("config")
	outputDir := cmd.String("output")

	// Load config toml file
	config := loadProject(projectFile)
	validateProjectConfig(config)
	config.OutputDir = outputDir

	return config, path.Dir(projectFile)
}

func validateProjectConfig(config interfaces.ProjectConfig) {
	if config.Project.Entrypoint == "" {
		log.Fatalln("Entrypoint not specified")
	}
	if config.Project.Name == "" {
		log.Fatalln("Name not specified")
	}
}
