package cmd

import (
	"fmt"

	"github.com/antham/yogo/v3/internal/client"
	"github.com/spf13/cobra"
)

var inboxDeleteCmd = &cobra.Command{
	Use:   "delete <inbox> <offset>",
	Short: "Delete email at given position in inbox",
	RunE:  inboxDelete(newInbox[client.MailHTMLDoc]),
	Args:  cobra.ExactArgs(2),
}

func inboxDelete(inboxBuilder inboxBuilder) cobraCmd {
	return func(cmd *cobra.Command, args []string) error {
		identifier := normalizeInboxName(args[0])
		offset, err := parseOffset(args[1])
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
		if err := in.Delete(offset - 1); err != nil {
			return err
		}
		cmd.Println(success(fmt.Sprintf(`Email "%d" successfully deleted`, offset)))
		return nil
	}
}

func init() {
	inboxCmd.AddCommand(inboxDeleteCmd)
}
