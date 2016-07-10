package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

// inboxShowCmd show full email
var inboxShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show full email",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			perror(fmt.Errorf("Two arguments mandatory"))

			errorExit()
		}

		offset, err := strconv.Atoi(args[1])

		if err != nil {
			perror(fmt.Errorf(`"%s" must be an integer`, args[1]))

			errorExit()
		}

		in, err := inbox.ParseInboxPages(args[0], offset)

		if err != nil {
			perror(err)

			errorExit()
		}

		in.Parse(offset-1)
		mail := in.Get(offset-1)
		renderMail(mail)
	},
}

func init() {
	inboxCmd.AddCommand(inboxShowCmd)
}
