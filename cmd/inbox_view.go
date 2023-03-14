package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/fatih/color"

	"github.com/antham/yogo/inbox"
)

var ErrSomethingWrongOccurred = errors.New("something wrong occurred")

func computeInboxMailOutput(in Inbox) (string, error) {
	JSON, err := computeJSONOutput(in)
	if err != nil {
		return "", err
	}
	if JSON != nil {
		return *JSON, nil
	}

	if in.Count() == 0 {
		return "", errors.New("Inbox is empty")
	}

	output := ""
	for index, mail := range in.GetMails() {
		var spam string
		if mail.IsSPAM {
			spam = " [SPAM]"
		}
		output = output + fmt.Sprintf(" %s %s%s%s\n", color.GreenString(fmt.Sprintf("%d", index+1)), color.YellowString(mail.Sender.Mail), color.YellowString(mail.Sender.Name), color.RedString(spam))
		output = output + fmt.Sprintf(" %s\n\n", color.CyanString(mail.Title))
	}
	return output, nil
}

func computeMailOutput(mail *inbox.Mail) (string, error) {
	JSON, err := computeJSONOutput(*mail)
	if err != nil {
		return "", err
	}
	if JSON != nil {
		return *JSON, nil
	}

	output := "---\n"
	if mail.Sender.Name == "" {
		output = output + fmt.Sprintf("From  : %s\n", color.MagentaString(mail.Sender.Mail))
	} else {
		output = output + fmt.Sprintf("From  : %s <%s>\n", color.MagentaString(mail.Sender.Name), color.MagentaString(mail.Sender.Mail))
	}
	output = output + fmt.Sprintf("Title : %s\n", color.YellowString(mail.Title))
	output = output + fmt.Sprintf("Date  : %s\n", color.GreenString(mail.Date.Format("2006-01-02 15:04")))
	output = output + "---\n"
	output = output + color.CyanString(mail.Body)
	output = output + "\n---\n"
	return output, nil
}

func computeJSONOutput(d interface{}) (*string, error) {
	if dumpJSON {
		data, err := json.Marshal(d)
		if err != nil {
			return nil, ErrSomethingWrongOccurred
		}
		s := string(data)
		return &s, nil
	}
	return nil, nil
}
