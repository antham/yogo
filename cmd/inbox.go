package cmd

import (
	"github.com/spf13/cobra"
)

type inboxBuilder func(string) (Inbox, error)

// inboxCmd represents the inbox command
var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "Handle inbox messages",
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			perror(err)
			errorExit()
		}
	},
}

func init() {
	RootCmd.AddCommand(inboxCmd)
}
