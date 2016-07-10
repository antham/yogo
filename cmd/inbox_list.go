package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

// inboxListCmd get all emails in an inbox
var inboxListCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all emails in an inbox",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			perror(fmt.Errorf("Two arguments mandatory"))

			errorExit()
		}

		limit, err := strconv.Atoi(args[1])

		if err != nil {
			perror(fmt.Errorf(`"%s" must be an integer`, args[1]))

			errorExit()
		}

		in, err := inbox.ParseInboxPages(args[0], limit)

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
