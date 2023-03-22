package cmd

import (
	"github.com/spf13/cobra"
)

var inboxListCmd = &cobra.Command{
	Use:   "list <inbox> <offset>",
	Short: "Get all emails from an inbox",
	RunE:  inboxList(newInbox),
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
		mail, err := computeInboxMailOutput(in, dumpJSON)
		if err != nil {
			return err
		}
		cmd.Println(mail)
		return nil
	}
}

func init() {
	inboxCmd.AddCommand(inboxListCmd)
}
