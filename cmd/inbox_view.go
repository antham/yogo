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
	if dumpJSON {
		data, err := json.Marshal(*in)
		if err != nil {
			perror(ErrSomethingWrongOccurred)
			errorExit()
		}

		output(string(data))
		successExit()
	}

	if in.Count() == 0 {
		info("Inbox is empty")

		successExit()
	}

	for index, mail := range in.GetAll() {
		output(fmt.Sprintf(" %s %s\n", color.GreenString(fmt.Sprintf("%d", index+1)), color.YellowString(mail.Title)))
		output(fmt.Sprintf(" %s\n\n", color.CyanString(*mail.SumUp)))
	}
}

func renderMail(mail *inbox.Mail) {
	if dumpJSON {
		data, err := json.Marshal(*mail)
		if err != nil {
			perror(ErrSomethingWrongOccurred)
			errorExit()
		}

		output(string(data))
		successExit()
	}

	output("---\n")
	if mail.Sender.Name == "" {
		output(fmt.Sprintf("From  : <%s>\n", color.MagentaString(mail.Sender.Mail)))
	} else {
		output(fmt.Sprintf("From  : %s <%s>\n", color.MagentaString(mail.Sender.Name), color.MagentaString(mail.Sender.Mail)))
	}
	output(fmt.Sprintf("Title : %s\n", color.YellowString(mail.Title)))
	output(fmt.Sprintf("Date  : %s\n", color.BlueString(mail.Date.Format("2006-01-02 15:04"))))
	output("---\n")
	output(color.CyanString(mail.Body))
	output("\n---\n")
}
