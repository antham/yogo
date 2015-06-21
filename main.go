package main

import "gopkg.in/alecthomas/kingpin.v2"
import "os"
import "github.com/antham/yogo/mailbox"
import "github.com/antham/yogo/view"

var (
	app     = kingpin.New("yogo", "Interact with yopmail from command line")
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	mailboxArgs      = app.Command("mailbox", "Handle mail")
	mailboxLimitArgs = mailboxArgs.Flag("limit", "Maximal number of messages to fetch").Default("1").Int()
	mailboxMailArgs  = mailboxArgs.Arg("mail", "mail").Required().String()
	mailboxFlushArgs = mailboxArgs.Flag("flush", "Flush inbox").Bool()

	mailArgs         = app.Command("mail", "Handle mail")
	mailMailArgs     = mailArgs.Arg("mail", "mail").Required().String()
	mailPositionArgs = mailArgs.Arg("position", "Position in mailbox").Default("1").Int()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case mailboxArgs.FullCommand():
		mailbox := mailbox.NewMailbox(*mailboxMailArgs)
		mails := mailbox.GetMails(*mailboxLimitArgs)
		view.OutputMails(mails)
	case mailArgs.FullCommand():
		mailbox := mailbox.NewMailbox(*mailMailArgs)
		mails := mailbox.GetMails(*mailPositionArgs)
		mail := mailbox.GetMail(mails[*mailPositionArgs-1].Id)
		view.OutputCompleteMail(mail)
	}
}
