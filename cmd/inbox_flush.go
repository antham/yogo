package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

// inboxFlushCmd flush all emails in an inbox
var inboxFlushCmd = &cobra.Command{
	Use:   "flush",
	Short: "Flush all emails in an inbox",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			perror(fmt.Errorf("One argument mandatory"))

			errorExit()
		}

		in, err := inbox.ParseInboxPages(parseMailAndOffsetArgs([]string{args[0], "1"}))
		if err != nil {
			perror(err)
			errorExit()
		}

		if err := in.Flush(); err != nil {
			perror(err)
			errorExit()
		}

		success(fmt.Sprintf(`Inbox "%s" successfully flushed`, args[0]))
	},
}

func init() {
	inboxCmd.AddCommand(inboxFlushCmd)
}
