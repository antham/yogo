package view

import "github.com/antham/yogo/mailbox"
import "github.com/antham/yogo/utils"
import "fmt"

func OutputMails(mails []*mailbox.Mail) {
	for _, mail := range mails {
		OutputMail(mail)
	}
}

func OutputMail(mail *mailbox.Mail) {
	fmt.Printf("%s\n%s\n", utils.Yellow(mail.Title), utils.Cyan(mail.SumUp))
	fmt.Printf("---\n")
}

func OutputCompleteMail(mail *mailbox.Mail) {

	fmt.Printf("---\n")
	fmt.Printf("From  : %s <%s>\n", utils.Magenta(mail.FromString), utils.Magenta(mail.FromMail))
	fmt.Printf("Title : %s\n", utils.Yellow(mail.Title))
	fmt.Printf("Date  : %s\n", utils.Blue(mail.Date.Format("2006-01-02 15:04")))
	fmt.Printf("---\n")
	fmt.Printf(utils.Cyan(mail.Body))
	fmt.Printf("\n---\n")
}
