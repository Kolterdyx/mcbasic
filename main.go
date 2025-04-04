package main

import (
	"context"
	"embed"
	"github.com/Kolterdyx/mcbasic/internal/command"
	"github.com/Kolterdyx/mcbasic/internal/interfaces"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
	"os"
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
		Authors: []any{
			(interfaces.Author{
				Name:  "Ciro García Belmonte (Kolterdyx)",
				Email: "info@cirogarcia.dev",
			}).String(),
		},
		Copyright: "(c) 2025 Ciro García Belmonte",
		Usage:     "A Minecraft datapack compiler",
		UsageText: "mcbasic [global options] <command> [command options]",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "enable debug mode",
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
					return command.Build(cmd, builtinHeaders, libs)
				},
			},
			{
				Name:  "init",
				Usage: "Initialize a new project",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Value: "MyProject",
						Usage: "Name of the project",
					},
					&cli.StringFlag{
						Name:  "namespace",
						Value: "myproject",
						Usage: "Namespace of the project",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return command.Init(cmd)
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cli.ShowAppHelpAndExit(cmd, 0)
			return nil
		},
	}
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
