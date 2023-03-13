package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

var inboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all emails in an inbox",
	RunE: inboxList(
		func(name string) (Inbox, error) {
			in, err := inbox.NewInbox(name)
			return Inbox(in), err
		},
	),
}

func inboxList(inboxBuilder inboxBuilder) cobraCmd {
	return func(cmd *cobra.Command, args []string) error {
		identifier, offset := parseMailAndOffsetArgs(args)

		in, err := inboxBuilder(identifier)
		if err != nil {
			return err
		}
		if err := in.ParseInboxPages(offset); err != nil {
			return err
		}
		return renderInboxMail(in)
	}
}

func init() {
	inboxCmd.AddCommand(inboxListCmd)
}
