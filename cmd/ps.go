package cmd

import (
	"engine/internal/containers"
	"engine/internal/pkg/util"
	"fmt"
	"io"

	"github.com/urfave/cli"
)

func PsCommand() cli.Command {
	return cli.Command{
		Name:  "ps",
		Usage: "List containers",
		Action: func(c *cli.Context) error {
			infos, err := containers.UpdateAllContainerStatus()
			if err != nil {
				return err
			}
			fmt.Fprintln(c.App.Writer, "CONTAINER-ID\tIMAGE\tCOMMAND\tSTATUS\tCREATEDAT")
			for id, info := range infos {
				fmt.Fprintf(c.App.Writer, "%s\t%s\t%s\t", id, info.Image, info.Command)
				printColoredStatus(c.App.Writer, info.Status) // 彩色状态
				fmt.Fprintf(c.App.Writer, "\t%s\n", util.FormatTimeAgo(info.Created))
			}
			return nil
		},
	}
}

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func printColoredStatus(w io.Writer, status string) {
	switch status {
	case "exited":
		fmt.Fprintf(w, "%s%s%s", colorRed, status, colorReset)
	case "running":
		fmt.Fprintf(w, "%s%s%s", colorGreen, status, colorReset)
	default:
		fmt.Fprint(w, status)
	}
}
