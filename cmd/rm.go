package cmd

import (
	"engine/internal/containers"
	"fmt"

	"github.com/urfave/cli"
)

func RmCommand() cli.Command {
	return cli.Command{
		Name:      "rm",
		Usage:     "Remove one or more containers",
		ArgsUsage: "CONTAINER [CONTAINER...]",

		Action: func(c *cli.Context) error {
			id := c.Args().First()
			err := containers.RemoveContainerInfo(id)
			if err != nil {
				fmt.Println("Error:", err)
				return err
			}
			fmt.Println("Remove container", id)
			return nil
		},
	}
}
