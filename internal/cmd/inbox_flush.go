package cmd

import (
	"fmt"

	"github.com/antham/yogo/internal/client"
	"github.com/spf13/cobra"
)

var inboxFlushCmd = &cobra.Command{
	Use:   "flush <inbox>",
	Short: "Flush all emails in an inbox",
	RunE:  inboxFlush(newInbox[client.MailHTMLDoc]),
	Args:  cobra.ExactArgs(1),
}

func inboxFlush(inboxBuilder inboxBuilder) cobraCmd {
	return func(cmd *cobra.Command, args []string) error {
		identifier := normalizeInboxName(args[0])
		in, err := inboxBuilder(identifier)
		if err != nil {
			return err
		}
		if err = in.ParseInboxPages(1); err != nil {
			return err
		}
		if err := in.Flush(); err != nil {
			return err
		}
		cmd.Println(success(fmt.Sprintf(`Inbox "%s" successfully flushed`, args[0])))
		return nil
	}
}

func init() {
	inboxCmd.AddCommand(inboxFlushCmd)
}
