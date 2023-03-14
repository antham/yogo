package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

var inboxShowCmd = &cobra.Command{
	Use:   "show <inbox> <offset>",
	Short: "Show full email at given position in inbox",
	RunE: inboxShow(
		func(name string) (Inbox, error) {
			in, err := inbox.NewInbox(name)
			return Inbox(in), err
		},
	),
	Args: cobra.ExactArgs(2),
}

func inboxShow(inboxBuilder inboxBuilder) cobraCmd {
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
		if err := checkOffset(in.Count(), offset); err != nil {
			return err
		}
		if err := in.Parse(offset - 1); err != nil {
			return err
		}

		mail := in.Get(offset - 1)
		if mail == nil {
			return nil
		}
		output, err := computeMailOutput(mail)
		if err != nil {
			return err
		}
		cmd.Println(output)
		return nil
	}
}

func init() {
	inboxCmd.AddCommand(inboxShowCmd)
}
