package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/antham/yogo/inbox"
)

var inboxFlushCmd = &cobra.Command{
	Use:   "flush <inbox>",
	Short: "Flush all emails in an inbox",
	RunE: inboxFlush(
		func(name string) (Inbox, error) {
			in, err := inbox.NewInbox(name)
			return Inbox(in), err
		},
	),
	Args: cobra.ExactArgs(1),
}

func inboxFlush(inboxBuilder inboxBuilder) cobraCmd {
	return func(cmd *cobra.Command, args []string) error {
		identifier, offset, err := parseMailAndOffsetArgs([]string{args[0], "1"})
		if err != nil {
			return err
		}

		in, err := inboxBuilder(identifier)
		if err != nil {
			return err
		}
		if err = in.ParseInboxPages(offset); err != nil {
			return err
		}
		if err := in.Flush(); err != nil {
			return err
		}

		success(fmt.Sprintf(`Inbox "%s" successfully flushed`, args[0]))
		return nil
	}
}

func init() {
	inboxCmd.AddCommand(inboxFlushCmd)
}
