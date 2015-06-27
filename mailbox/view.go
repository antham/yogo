package mailbox

import "github.com/antham/yogo/utils"
import "github.com/antham/yogo/mail"
import "fmt"
import "strings"
import "strconv"

func Render(mailbox *Mailbox) {
	mails := mailbox.GetAll()

	for index, mail := range mails {
		renderMail(mail, index, len(mails))
	}
}

func renderMail(mail *mail.Mail, index int, totalMails int) {
	totalSpaces := strings.Repeat(" ", len(strconv.Itoa(totalMails)))
	remainingSpaces := strings.Repeat(" ", len(strconv.Itoa(totalMails))-len(strconv.Itoa(index+1)))

	fmt.Printf(" %s%v %s\n", utils.Green(index+1), remainingSpaces, utils.Yellow(mail.Title))
	fmt.Printf(" %v %s\n\n", totalSpaces, utils.Cyan(mail.SumUp))
}

func RenderMessage(string string) {
	fmt.Println(utils.Cyan(string))
}
