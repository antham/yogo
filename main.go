package main

import "gopkg.in/alecthomas/kingpin.v2"
import "os"
import "github.com/antham/yogo/mailbox"
import "github.com/antham/yogo/view"

var (
	app     = kingpin.New("yogo", "Interact with yopmail from command line")
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	mailboxArgs      = app.Command("mailbox", "Handle mail")
	mailboxEmailArgs = mailboxArgs.Arg("email", "Email").Required().String()
	mailboxLimitArgs = mailboxArgs.Flag("limit", "Maximal number of messages to fetch").Int()
	mailboxFlushArgs = mailboxArgs.Flag("flush", "Flush inbox").Bool()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	mailbox := mailbox.NewMailbox(*mailboxEmailArgs)
	mails := mailbox.GetMails(*mailboxLimitArgs)
	view.OutputMails(mails)
}
