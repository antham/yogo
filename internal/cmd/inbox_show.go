package cmd

import (
	"github.com/antham/yogo/internal/client"
	"github.com/spf13/cobra"
)

var isSourceMail bool

var inboxShowCmd = &cobra.Command{
	Use:   "show <inbox> <offset>",
	Short: "Show full email at given position in inbox",
	RunE:  inboxShow(newInbox[client.MailHTMLDoc]),
	Args:  cobra.ExactArgs(2),
}

func inboxShow(inboxBuilder inboxBuilder) cobraCmd {
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

		mail, err := in.Fetch(offset - 1)
		if err != nil {
			return err
		}
		if mail == nil {
			return nil
		}

		var output string
		if dumpJSON {
			output, err = mail.JSON()
			if err != nil {
				return err
			}
		} else {
			output, err = mail.Coloured()
			if err != nil {
				return err
			}
		}

		cmd.Println(output)
		return nil
	}
}

func init() {
	inboxShowCmd.Flags().BoolVar(&isSourceMail, "source", false, "Get the source of the message if enabled")
	inboxCmd.AddCommand(inboxShowCmd)
}
