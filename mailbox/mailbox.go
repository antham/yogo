package mailbox

import "github.com/PuerkitoBio/goquery"
import "fmt"
import "log"
import "regexp"

var baseUrl = "http://www.yopmail.com/en/inbox.php?login=%v&p=r&d=&ctrl=&scrl=&spam=true&v=2.6&r_c=&id="

type MailBox struct {
	mail  string
	mails []string
}

func NewMailBox(mail string) *MailBox {
	return &MailBox{
		mail: mail,
	}
}

func (m *MailBox) GetMails(limit int) []*Mail {
	var mails []*Mail

	doc, err := goquery.NewDocument(fmt.Sprintf(baseUrl, m.mail))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.um").Each(func(i int, s *goquery.Selection) {
		re := regexp.MustCompile("mail.php.b=.*?id=(.*)")

		idUrl, _ := s.Find("a.lm").Attr("href")

		mail := &Mail{
			id: re.FindStringSubmatch(idUrl)[1],
			title: s.Find("span.lmf").Text(),
			sumUp: s.Find("span.lms").Text(),
		}

		mails = append(mails, mail)
	})

	return mails
}
