package cmd

import (
	"engine/internal/pkg/repository"
	"fmt"

	"github.com/urfave/cli"
)

func Images() cli.Command {

	return cli.Command{
		Name:        "images",
		Usage:       "List images",
		Description: "List all images in local repository",
		Action: func(c *cli.Context) error {
			list, err := repository.GetImagesList()
			if err != nil {
				return err
			}
			fmt.Println("IMAGE")
			for _, image := range list {
				fmt.Println(image)
			}
			return nil
		},
	}
}
