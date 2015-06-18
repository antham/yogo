package main

import "gopkg.in/alecthomas/kingpin.v2"
import "os"
import "github.com/antham/yogo/mailbox"

var (
	app     = kingpin.New("yogo", "Interact with yopmail from command line")
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()

	mailbox      = app.Command("mailbox", "Handle mail")
	mailboxEmail = mailbox.Arg("email", "Email").Required().String()
	mailboxMax   = mailbox.Flag("limit", "Maximal number of messages to fetch").Int()
	mailboxFlush = mailbox.Flag("flush", "Flush inbox").Bool()
)

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
