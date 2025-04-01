package cmd

import (
	"engine/internal/pkg/repository"
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
			for _, image := range list {
				println(image)
			}
			return nil
		},
	}
}
