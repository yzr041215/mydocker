package cmd

import "github.com/urfave/cli"

func StopCommand() cli.Command {
	return cli.Command{
		Name:  "stop",
		Usage: "Stop a running container,ag: mydocker stop CONTAINER",

		Action: func(c *cli.Context) error {
			containerID := c.Args().First()
			if containerID == "" {
				return cli.NewExitError("Please specify the container ID or name", 1)
			}

			return nil
		},
	}
}
