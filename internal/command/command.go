package command

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

type Data struct {
	Stdlib   embed.FS
	Embedded embed.FS
}

func DebugFlagHandler(cmd *cli.Command) {
	if cmd.Bool("debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
