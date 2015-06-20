package view

import "github.com/fatih/color"
import "github.com/antham/yogo/mailbox"
import "fmt"

var blue = color.New(color.FgBlue).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()

func OutputMails(mails []*mailbox.Mail) {
	for _,mail := range(mails) {
		OutputMail(mail)
	}
}

func OutputMail(mail *mailbox.Mail) {
	fmt.Printf("%s\n%s\n", yellow(mail.Title), cyan(mail.SumUp))
	fmt.Printf("---\n")
}
