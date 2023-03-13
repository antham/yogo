package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

var inboxShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show full email at given position in inbox",
	RunE: inboxShow(
		func(name string) (Inbox, error) {
			in, err := inbox.NewInbox(name)
			return Inbox(in), err
		},
	),
}

func inboxShow(inboxBuilder inboxBuilder) cobraCmd {
	return func(cmd *cobra.Command, args []string) error {
		identifier, offset := parseMailAndOffsetArgs(args)
		in, err := inboxBuilder(identifier)
		if err != nil {
			return err
		}
		if err := in.ParseInboxPages(offset); err != nil {
			return err
		}
		if err := checkOffset(in.Count(), offset); err != nil {
			return err
		}
		if err := in.Parse(offset - 1); err != nil {
			return err
		}

		mail := in.Get(offset - 1)
		if mail != nil {
			renderMail(mail)
		}
		return nil
	}
}

func init() {
	inboxCmd.AddCommand(inboxShowCmd)
}
