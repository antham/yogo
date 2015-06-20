package mailbox

import "github.com/PuerkitoBio/goquery"
import "fmt"
import "log"
import "regexp"


var baseUrl = "http://www.yopmail.com/en/inbox.php?login=%v&p=%v&d=&ctrl=&scrl=&spam=true&v=2.6&r_c=&id="
var mailPerPage = 15

type Mailbox struct {
	mail  string
}

func NewMailbox(mail string) *Mailbox {
	return &Mailbox{
		mail: mail,
	}
}

func (m *Mailbox) GetMails(limit int) []*Mail {
	var mails []*Mail

	for counter := 1; counter <= int(limit / mailPerPage) + 1; counter++ {

		doc, err := goquery.NewDocument(fmt.Sprintf(baseUrl, m.mail, counter))
		if err != nil {
			log.Fatal(err)
		}

		doc.Find("div.um").Each(func(i int, s *goquery.Selection) {
			re := regexp.MustCompile("mail.php.b=.*?id=(.*)")

			idUrl, _ := s.Find("a.lm").Attr("href")

			matches := re.FindStringSubmatch(idUrl)

			if len(matches) == 2 {
				mail := &Mail{
					id: matches[1],
					Title: s.Find("span.lmf").Text(),
					SumUp: s.Find("span.lms").Text(),
				}

				mails = append(mails, mail)
			}
		})
	}

	if limit >= len(mails) {
		return mails
	}

	return mails[:limit]
}
