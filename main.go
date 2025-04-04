package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
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

//go:embed version.txt
var version string

//go:embed libs
var libs embed.FS

//go:embed headers
var builtinHeaders embed.FS

func main() {

	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})

	cmd := &cli.Command{
		Name:    "mcbasic",
		Version: version,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Value:   false,
				Usage:   "Print version",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "Enable debug mode",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "Compile the project into a datapack",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "config",
						Value: "project.toml",
						Usage: "Path to the project config file",
					},
					&cli.StringFlag{
						Name:  "output",
						Value: "build",
						Usage: "Output directory. The resulting datapack will be inside this directory as <output>/<project_name>",
					},
					&cli.BoolFlag{
						Name:  "enable-traces",
						Value: false,
						Usage: "Enable traces",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {

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
						return fmt.Errorf("scanner had an error")
					}
					log.Debug("Tokens scanned successfully")

					headers, err := loadHeaders(config.Dependencies.Headers, projectRoot)
					if err != nil {
						return err
					}

					parser_ := parser.Parser{Tokens: tokens, Headers: headers}
					program := parser_.Parse()
					if parser_.HadError {
						return fmt.Errorf("parser had an error")
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
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
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
