package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

// inboxListCmd get all emails in an inbox
var inboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all emails in an inbox",
	Run: func(cmd *cobra.Command, args []string) {
		identifier, offset := parseMailAndOffsetArgs(args)

		in, err := inbox.ParseInboxPages(identifier, offset)

		if err != nil {
			perror(err)

			errorExit()
		}

		renderInboxMail(in)
	},
}

func init() {
	inboxCmd.AddCommand(inboxListCmd)
}
