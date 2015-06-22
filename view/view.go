package view

import "github.com/fatih/color"
import "github.com/antham/yogo/mailbox"
import "fmt"

var blue = color.New(color.FgBlue).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var magenta = color.New(color.FgMagenta).SprintFunc()

func OutputMails(mails []*mailbox.Mail) {
	for _, mail := range mails {
		OutputMail(mail)
	}
}

func OutputMail(mail *mailbox.Mail) {
	fmt.Printf("%s\n%s\n", yellow(mail.Title), cyan(mail.SumUp))
	fmt.Printf("---\n")
}

func OutputCompleteMail(mail *mailbox.Mail) {

	fmt.Printf("---\n")
	fmt.Printf("From  : %s <%s>\n", magenta(mail.FromString), magenta(mail.FromMail))
	fmt.Printf("Title : %s\n", yellow(mail.Title))
	fmt.Printf("Date  : %s\n", blue(mail.Date.Format("2006-01-02 15:04")))
	fmt.Printf("---\n")
	fmt.Printf(cyan(mail.Body))
	fmt.Printf("\n---\n")
}
