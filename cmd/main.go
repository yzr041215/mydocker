package cmd

import (
	"os"

	"github.com/urfave/cli"
)

func Cmd() {
	rootCmd := cli.NewApp()
	rootCmd.Name = "MyDocker"
	rootCmd.Usage = "A simple Docker client"
	rootCmd.Version = "1.0.0"
	rootCmd.Commands = append(rootCmd.Commands, runCommand)

	rootCmd.Commands = append(rootCmd.Commands, runCommand, Images(), PullCommand())
	rootCmd.Commands = append(rootCmd.Commands, PsCommand(), RmCommand(), LogsCommand(), NetworkCommand())
	err := rootCmd.Run(os.Args)
	if err != nil {
		return
	}
}
