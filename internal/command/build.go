package command

import (
	"context"
	"embed"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/Kolterdyx/mcbasic/internal/compiler"
	frontend "github.com/Kolterdyx/mcbasic/internal/frontend"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"os"
	"path"
	"strings"
)

var BuildCommand = &cli.Command{
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
		data := ctx.Value("data").(*Data)
		return Build(cmd, data.BuiltinHeaders, data.Libs)
	},
}

func Build(cmd *cli.Command, builtinHeaders, libs embed.FS) error {
	DebugFlagHandler(cmd)

	config, projectRoot := parseArgs(cmd)

	entrypoint := config.Project.Entrypoint

	front := frontend.NewFrontend(projectRoot)

	err := front.Parse(entrypoint)
	if err != nil {
		return err
	}
	// Remove the contents of the output directory
	err = os.RemoveAll(path.Join(config.OutputDir, config.Project.Name))
	if err != nil {
		return err
	}

	c := compiler.NewCompiler(config, headers, libs)
	err = c.Compile(program)
	if err != nil {
		log.Error(err)
		return cli.Exit("There was an error compiling the program", 1)
	}
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
		headerFile, err := builtinHeaders.ReadFile(headerPath)
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
