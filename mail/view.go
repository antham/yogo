package mail

import "github.com/antham/yogo/utils"
import "fmt"

func Render(mail *Mail) {

	fmt.Printf("---\n")
	fmt.Printf("From  : %s <%s>\n", utils.Magenta(mail.FromString), utils.Magenta(mail.FromMail))
	fmt.Printf("Title : %s\n", utils.Yellow(mail.Title))
	fmt.Printf("Date  : %s\n", utils.Blue(mail.Date.Format("2006-01-02 15:04")))
	fmt.Printf("---\n")
	fmt.Printf(utils.Cyan(mail.Body))
	fmt.Printf("\n---\n")
}
