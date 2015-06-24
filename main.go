package main

import "gopkg.in/alecthomas/kingpin.v2"
import "os"
import mailboxmod "github.com/antham/yogo/mailbox"
import mailmod "github.com/antham/yogo/mail"

var (
	app     = kingpin.New("yogo", "Interact with yopmail from command line")
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	mailboxArgs       = app.Command("mailbox", "Manage mailbox")
	mailboxLimitArgs  = mailboxArgs.Flag("limit", "Maximal number of messages to fetch").Default("1").Int()
	mailboxMailArgs   = mailboxArgs.Arg("mail", "Targeted inbox").Required().String()
	mailboxActionArgs = mailboxArgs.Arg("action", "Action to do").Default("list").Enum("list", "flush")

	mailArgs         = app.Command("mail", "Manage mail")
	mailMailArgs     = mailArgs.Arg("mail", "Targeted inbox").Required().String()
	mailPositionArgs = mailArgs.Arg("position", "Position in mailbox").Default("1").Int()
	mailActionArgs   = mailArgs.Arg("action", "Action to do").Default("read").Enum("read", "delete")
)

func main() {
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case mailboxArgs.FullCommand():
		callMailboxAction(mailboxActionArgs)
	case mailArgs.FullCommand():
		callMailAction(mailActionArgs)
	}
}

func callMailboxAction(action *string) {
	mailbox := mailboxmod.NewMailbox(*mailboxMailArgs)

	switch *action {
	case "list":
		mailbox.Fetch(*mailboxLimitArgs)
		if mailbox.Count() != 0 {
			mailboxmod.Render(mailbox)
		} else {
			mailboxmod.RenderMessage("Mo mails found")
		}
	case "flush":
		mailbox.Flush()
	}
}

func callMailAction(action *string) {
	mailbox := mailboxmod.NewMailbox(*mailMailArgs)
	mailbox.Fetch(*mailPositionArgs)
	mail := mailbox.Get(*mailPositionArgs - 1)

	if mail != nil {
		switch *action {
		case "read":
			mail.Fetch()
			mailmod.Render(mail)
		}
	} else {
		mailmod.RenderMessage("No mail found")
	}
}
