package cmd

import (
	"github.com/spf13/cobra"
)

type inboxBuilder func(string) (Inbox, error)

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "Handle inbox messages",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(inboxCmd)
}
