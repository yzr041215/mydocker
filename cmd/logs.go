package cmd

import (
	"engine/internal/containers"
	"fmt"

	"github.com/urfave/cli"
)

func LogsCommand() cli.Command {
	return cli.Command{
		Name:      "logs",
		Usage:     "Show logs of a container",
		ArgsUsage: "CONTAINER", // TODO: Update usage message
		Action: func(c *cli.Context) error {
			id := c.Args().First()
			if id == "" {
				fmt.Println("Container name or ID is required")
				return cli.NewExitError("Container name or ID is required", 1)
			}
			fmt.Println("Logs for container ID: ", id, " :")
			logs, err := containers.GetLog(id)
			if err != nil {
				return cli.NewExitError(err, 1)
			}
			fmt.Println(logs)
			return nil
		},
	}
}
