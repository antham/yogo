package cmd

import (
	"github.com/antham/yogo/v4/internal/client"
	"github.com/antham/yogo/v4/internal/inbox"
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

func newInbox[M client.MailDoc](name string) (Inbox, error) {
	in, err := inbox.NewInbox[M](name)
	return Inbox(in), err
}
