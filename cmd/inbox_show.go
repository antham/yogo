package cmd

import (
	"github.com/spf13/cobra"
)

var inboxShowCmd = &cobra.Command{
	Use:   "show <inbox> <offset>",
	Short: "Show full email at given position in inbox",
	RunE:  inboxShow(newInbox),
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
		output, err := computeMailOutput(mail, dumpJSON)
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
