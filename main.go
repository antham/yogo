package main

import "gopkg.in/alecthomas/kingpin.v2"
import "os"
import mailboxmod "github.com/antham/yogo/mailbox"
import mailmod "github.com/antham/yogo/mail"

var (
	app     = kingpin.New("yogo", "Interact with yopmail from command line")
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	mailboxArgs      = app.Command("mailbox", "Manage mailbox")
	mailboxLimitArgs = mailboxArgs.Flag("limit", "Maximal number of messages to fetch").Default("1").Int()
	mailboxFlushArgs = mailboxArgs.Flag("flush", "Flush inbox").Bool()
	mailboxMailArgs  = mailboxArgs.Arg("mail", "Targeted inbox").Required().String()

	mailArgs         = app.Command("mail", "Manage mail")
	mailMailArgs     = mailArgs.Arg("mail", "Targeted inbox").Required().String()
	mailPositionArgs = mailArgs.Arg("position", "Position in mailbox").Default("1").Int()
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case mailboxArgs.FullCommand():
		mailbox := mailboxmod.NewMailbox(*mailboxMailArgs)
		mails := mailbox.GetMails(*mailboxLimitArgs)
		mailboxmod.OutputMails(mails)
	case mailArgs.FullCommand():
		mailbox := mailboxmod.NewMailbox(*mailMailArgs)
		mails := mailbox.GetMails(*mailPositionArgs)
		mail := mailbox.GetMail(mails[*mailPositionArgs-1].Id)
		mailmod.OutputMail(mail)
	}
}
