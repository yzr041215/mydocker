package cmd

import (
	"engine/internal/pkg/registry"
	"fmt"

	"github.com/urfave/cli"
)

// go run main.go  pull -image mysql
func PullCommand() cli.Command {

	return cli.Command{
		Name:        "pull",
		Usage:       "Pull images",
		Description: "Pull images from remote repository",
		Action: func(c *cli.Context) error {

			if err := registry.Pull(c.Args().First()); err != nil {
				fmt.Println(err)
				return err
			}

			return nil
		},
	}
}
