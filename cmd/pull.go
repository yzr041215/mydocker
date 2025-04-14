package cmd

import (
	"engine/internal/pkg/registry"

	"github.com/urfave/cli"
)

// go run main.go  pull -image mysql
func PullCommand() cli.Command {

	return cli.Command{
		Name:        "pull",
		Usage:       "Pull images",
		Description: "Pull images from remote repository",
		Action: func(c *cli.Context) error {

			registry.Pull(c.Args().First())
			return nil
		},
	}
}
