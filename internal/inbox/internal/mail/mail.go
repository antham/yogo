package mail

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/antham/yogo/internal/client"
)

const noDataToDisplayMsg = "[no data to display]"

type Identifier interface {
	SetID(string)
}

type Render interface {
	Coloured() (string, error)
	JSON() (string, error)
}

type RenderIdentifier interface {
	Identifier
	Render
}

func Parse[M client.MailDoc](doc M) (RenderIdentifier, error) {
	var m RenderIdentifier
	switch any(doc).(type) {
	case client.MailHTMLDoc:
		mail := &HTMLMail{}
		doc.Find("body div.fl .ellipsis").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 0:
				mail.Title = strings.TrimSpace(s.Text())
			case 1:
				mail.Sender = &Sender{}
				mail.Sender.Name, mail.Sender.Mail = parseFrom(s.Text())
			case 2:
				mail.Date = parseDate(strings.Join(strings.Fields(s.Text()), " "))
			}
		})
		mail.Body = parseHTML(doc.Find("div#mail").Html())
		m = mail
	}
	return m, nil
}
