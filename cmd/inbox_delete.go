package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

// inboxDeleteCmd delete an email in inbox
var inboxDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete email at given position in inbox",
	Run: func(cmd *cobra.Command, args []string) {

		identifier, offset := parseMailAndOffsetArgs(args)

		in, err := inbox.ParseInboxPages(identifier, offset)

		if err != nil {
			perror(err)
			errorExit()
		}

		checkOffset(in.Count(), offset)

		if err := in.Delete(offset - 1); err != nil {
			perror(err)
			errorExit()
		}
		success(fmt.Sprintf(`Email "%d" successfully deleted`, offset))
	},
}

func init() {
	inboxCmd.AddCommand(inboxDeleteCmd)
}
