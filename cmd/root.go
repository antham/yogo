package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type cobraCmd func(*cobra.Command, []string) error

var dumpJSON = false

var RootCmd = &cobra.Command{
	Use:   "yogo",
	Short: "Interact with yopmail from command-line",
	Long:  `Check yopmail mails and inboxes from command line.`,
}

func Execute() {
	RootCmd.PersistentFlags().BoolVar(&dumpJSON, "json", false, "Dump the output as json")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
