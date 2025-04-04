package command

import "github.com/urfave/cli/v3"

func Init(cmd *cli.Command) error {

	name := cmd.String("name")
	if name == "" {

	}

	return nil
}
