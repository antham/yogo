package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fatih/color"

	"github.com/antham/yogo/inbox"
)

var ErrSomethingWrongOccurred = errors.New("something wrong occurred")

func renderInboxMail(in Inbox) error {
	renderJSON(in)

	if in.Count() == 0 {
		info("Inbox is empty")
		return nil
	}

	for index, mail := range in.GetMails() {
		var spam string
		if mail.IsSPAM {
			spam = " [SPAM]"
		}

		output(fmt.Sprintf(" %s %s%s%s\n", color.GreenString(fmt.Sprintf("%d", index+1)), color.YellowString(mail.Sender.Mail), color.YellowString(mail.Sender.Name), color.RedString(spam)))
		output(fmt.Sprintf(" %s\n\n", color.CyanString(mail.Title)))
	}
	return nil
}

func renderMail(mail *inbox.Mail) {
	renderJSON(*mail)

	const noDataToDisplay = "No data to display"

	output("---\n")
	switch {
	case mail.Sender.Name == "" && mail.Sender.Mail == "":
		output(fmt.Sprintf("From  : %s\n", color.RedString(noDataToDisplay)))
	case mail.Sender.Name == "":
		output(fmt.Sprintf("From  : %s\n", color.MagentaString(mail.Sender.Mail)))
	default:
		output(fmt.Sprintf("From  : %s <%s>\n", color.MagentaString(mail.Sender.Name), color.MagentaString(mail.Sender.Mail)))
	}
	output(fmt.Sprintf("Title : %s\n", color.YellowString(mail.Title)))

	if mail.Date == nil {
		output(fmt.Sprintf("Date  : %s\n", color.RedString(noDataToDisplay)))
	} else {
		output(fmt.Sprintf("Date  : %s\n", color.GreenString(mail.Date.Format("2006-01-02 15:04"))))
	}
	output("---\n")
	if mail.Body == "" {
		output(color.RedString(noDataToDisplay))
	} else {
		output(color.CyanString(mail.Body))
	}
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
