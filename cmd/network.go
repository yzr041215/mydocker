package cmd

import (
	"engine/internal/network"
	"fmt"

	"github.com/urfave/cli"
)

func NetworkCommand() cli.Command {
	return networkCommand
}

var networkCommand = cli.Command{
	Name:  "network",
	Usage: "container network commands",
	Subcommands: []cli.Command{
		{
			// mydeocker network create --driver bridge --subnet 192.168.0.0/24 my-network
			Name:  "create",
			Usage: "create a container network",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "driver",
					Usage: "network driver",
				},
				cli.StringFlag{
					Name:  "subnet",
					Usage: "subnet cidr",
				},
			},
			// mydeocker network create --driver bridge --subnet 192.168.0.0/16 mydocker-network
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("missing network name")
				}
				driver := context.String("driver")
				subnet := context.String("subnet")
				name := context.Args()[0]
				fmt.Println("create network", driver, subnet, name)
				err := network.CreateNetwork(driver, subnet, name)
				if err != nil {
					return fmt.Errorf("create network error: %+v", err)
				}
				return nil
			},
		},
		{
			// mydeocker network list
			Name:  "list",
			Usage: "list container network",
			Action: func(context *cli.Context) error {
				network.ListNetwork()
				return nil
			},
		},
		{
			Name:  "remove",
			Usage: "remove container network",
			Action: func(context *cli.Context) error {
				if len(context.Args()) < 1 {
					return fmt.Errorf("missing network name")
				}
				err := network.DeleteNetwork(context.Args()[0])
				if err != nil {
					return fmt.Errorf("remove network error: %+v", err)
				}
				return nil
			},
		},
	},
}
