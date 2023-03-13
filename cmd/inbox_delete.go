package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

// inboxDeleteCmd delete an email in inbox
var inboxDeleteCmd = &cobra.Command{
	Use:   "delete <inbox> <offset>",
	Short: "Delete email at given position in inbox",
	RunE: inboxDelete(
		func(name string) (Inbox, error) {
			in, err := inbox.NewInbox(name)
			return Inbox(in), err
		},
	),
	Args: cobra.ExactArgs(2),
}

func inboxDelete(inboxBuilder inboxBuilder) cobraCmd {
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
		if err := in.Delete(offset - 1); err != nil {
			return err
		}
		success(fmt.Sprintf(`Email "%d" successfully deleted`, offset))
		return nil
	}
}

func init() {
	inboxCmd.AddCommand(inboxDeleteCmd)
}
