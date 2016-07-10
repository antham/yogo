package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

// inboxShowCmd show full email
var inboxShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show full email at given position in inbox",
	Run: func(cmd *cobra.Command, args []string) {
		identifier, offset := parseMailAndOffsetArgs(args)

		in, err := inbox.ParseInboxPages(identifier, offset)

		if err != nil {
			perror(err)

			errorExit()
		}

		checkOffset(in.Count(), offset)

		in.Parse(offset - 1)
		mail := in.Get(offset - 1)
		renderMail(mail)
	},
}

func init() {
	inboxCmd.AddCommand(inboxShowCmd)
}
