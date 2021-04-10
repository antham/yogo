package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fatih/color"

	"github.com/antham/yogo/inbox"
)

var ErrSomethingWrongOccurred = errors.New("something wrong occurred")

func renderInboxMail(in *inbox.Inbox) {
	renderJSON(*in)

	if in.Count() == 0 {
		info("Inbox is empty")

		successExit()
	}

	for index, mail := range in.Mails {
		var spam string
		if mail.IsSPAM {
			spam = " [SPAM]"
		}

		output(fmt.Sprintf(" %s %s%s%s\n", color.GreenString(fmt.Sprintf("%d", index+1)), color.YellowString(mail.Sender.Mail), color.YellowString(mail.Sender.Name), color.RedString(spam)))
		output(fmt.Sprintf(" %s\n\n", color.CyanString(mail.Title)))
	}
}

func renderMail(mail *inbox.Mail) {
	renderJSON(*mail)

	output("---\n")
	if mail.Sender.Name == "" {
		output(fmt.Sprintf("From  : %s\n", color.MagentaString(mail.Sender.Mail)))
	} else {
		output(fmt.Sprintf("From  : %s <%s>\n", color.MagentaString(mail.Sender.Name), color.MagentaString(mail.Sender.Mail)))
	}
	output(fmt.Sprintf("Title : %s\n", color.YellowString(mail.Title)))
	output(fmt.Sprintf("Date  : %s\n", color.GreenString(mail.Date.Format("2006-01-02 15:04"))))
	output("---\n")
	output(color.CyanString(mail.Body))
	output("\n---\n")
}

func renderJSON(d interface{}) {
	if dumpJSON {
		data, err := json.Marshal(d)
		if err != nil {
			perror(ErrSomethingWrongOccurred)
			errorExit()
		}

		output(string(data))
		successExit()
	}
}
