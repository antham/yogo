package mailbox

import "github.com/antham/yogo/utils"
import "github.com/antham/yogo/mail"
import "fmt"

func Render(mailbox *Mailbox) {
	for _, mail := range mailbox.GetAll() {
		renderMail(mail)
	}
}

func renderMail(mail *mail.Mail) {
	fmt.Printf("%s\n%s\n", utils.Yellow(mail.Title), utils.Cyan(mail.SumUp))
	fmt.Printf("---\n")
}
