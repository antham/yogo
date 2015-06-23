package mailbox

import "github.com/PuerkitoBio/goquery"
import mailmod "github.com/antham/yogo/mail"
import "fmt"
import "log"
import "regexp"

var mailboxBaseUrl = "http://www.yopmail.com/en/inbox.php?login=%v&p=%v&d=&ctrl=&scrl=&spam=true&v=2.6&r_c=&id="
var mailPerPage = 15

type Mailbox struct {
	mail  string
	mails []*mailmod.Mail
}

func NewMailbox(mail string) *Mailbox {
	return &Mailbox{
		mail: mail,
	}
}

func (m *Mailbox) Fetch(limit int) {
	var mails []*mailmod.Mail

	for counter := 1; counter <= int(limit/mailPerPage)+1; counter++ {

		doc, err := goquery.NewDocument(fmt.Sprintf(mailboxBaseUrl, m.mail, counter))
		if err != nil {
			log.Fatal(err)
		}

		doc.Find("div.um").Each(func(i int, s *goquery.Selection) {

			id := func(s *goquery.Selection) string {
				re := regexp.MustCompile("mail.php.b=.*?id=(.*)")

				idUrl, _ := s.Find("a.lm").Attr("href")

				matches := re.FindStringSubmatch(idUrl)

				if len(matches) == 2 {
					return matches[1]
				}

				return ""
			}(s)

			if id != "" {
				mail := mailmod.NewMail(m.mail, id, s.Find("span.lmf").Text(), s.Find("span.lms").Text())

				mails = append(mails, mail)
			}
		})

		if limit >= len(mails) {
			m.mails = mails
		}
	}

	m.mails = mails[:limit]
}

func (m *Mailbox) GetAll() []*mailmod.Mail {
	return m.mails
}

func (m *Mailbox) Get(position int) *mailmod.Mail {
	return m.mails[position]
}
