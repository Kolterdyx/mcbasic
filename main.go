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
			command.BuildCommand,
			command.InitCommand,
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			cli.ShowAppHelpAndExit(cmd, 0)
			return nil
		},
	}

	ctx := context.WithValue(context.Background(), "data", &command.Data{
		BuiltinHeaders: builtinHeaders,
		Libs:           libs,
	})

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
