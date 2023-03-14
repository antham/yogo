package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

var inboxListCmd = &cobra.Command{
	Use:   "list <inbox> <offset>",
	Short: "Get all emails in an inbox",
	RunE: inboxList(
		func(name string) (Inbox, error) {
			in, err := inbox.NewInbox(name)
			return Inbox(in), err
		},
	),
	Args: cobra.ExactArgs(2),
}

func inboxList(inboxBuilder inboxBuilder) cobraCmd {
	return func(cmd *cobra.Command, args []string) error {
		identifier, offset, err := parseMailAndOffsetArgs(args)
		if err != nil {
			return err
		}

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
