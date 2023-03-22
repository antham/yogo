package cmd

import (
	"github.com/antham/yogo/inbox"
	"github.com/spf13/cobra"
)

type inboxBuilder func(string) (Inbox, error)

var inboxCmd = &cobra.Command{
	Use:   "inbox",
	Short: "Handle inbox messages",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	RootCmd.AddCommand(inboxCmd)
}

func newInbox(name string) (Inbox, error) {
	in, err := inbox.NewInbox(name)
	return Inbox(in), err
}
