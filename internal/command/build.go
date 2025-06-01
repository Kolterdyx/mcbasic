package command

import (
	"context"
	"embed"
	"github.com/BurntSushi/toml"
	frontend "github.com/Kolterdyx/mcbasic/internal/frontend"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"os"
	"path"
	"path/filepath"
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
		return Build(cmd, data.Stdlib, data.Embedded)
	},
}

func Build(cmd *cli.Command, builtinHeaders, libs embed.FS) error {
	DebugFlagHandler(cmd)

	config, projectRoot := parseArgs(cmd)

	entrypoint := config.Project.Entrypoint
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return err
	}
	front := frontend.NewFrontend(config, absProjectRoot, builtinHeaders, libs)

	absEntrypoint, err := filepath.Abs(path.Join(projectRoot, entrypoint))
	if err != nil {
		return err
	}
	err = front.Parse(absEntrypoint)
	if err != nil {
		return err
	}
	// Remove the contents of the output directory
	err = os.RemoveAll(path.Join(config.OutputDir, config.Project.Name))
	if err != nil {
		return err
	}

	errs := front.Compile()
	if len(errs) > 0 {
		for _, err := range errs {
			if err != nil {
				log.Error(err)
			}
		}
		return nil
	}
	log.Info("Compilation complete")
	return nil
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
