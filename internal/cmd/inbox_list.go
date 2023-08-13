package cmd

import (
	"github.com/antham/yogo/internal/client"
	"github.com/spf13/cobra"
)

var inboxListCmd = &cobra.Command{
	Use:   "list <inbox> <offset>",
	Short: "Get all emails from an inbox",
	RunE:  inboxList(newInbox[client.MailHTMLDoc]),
	Args:  cobra.ExactArgs(2),
}

func inboxList(inboxBuilder inboxBuilder) cobraCmd {
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

		var output string
		if dumpJSON {
			output, err = in.JSON()
			if err != nil {
				return err
			}
		} else {
			output, err = in.Coloured()
			if err != nil {
				return err
			}
		}

		cmd.Println(output)
		return nil
	}
}

func init() {
	inboxCmd.AddCommand(inboxListCmd)
}
