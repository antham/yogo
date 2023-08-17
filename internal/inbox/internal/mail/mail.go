package mail

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/antham/yogo/v4/internal/client"
	"io"
	gomail "net/mail"
	"strings"
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
				mail.Subject = strings.TrimSpace(s.Text())
			case 1:
				mail.Sender = &Sender{}
				mail.Sender.Name, mail.Sender.Mail = parseFrom(s.Text())
			case 2:
				mail.Date = parseDate(strings.Join(strings.Fields(s.Text()), " "))
			}
		})
		mail.Body = parseHTML(doc.Find("div#mail").Html())
		m = mail
	case client.MailSourceDoc:
		msg, err := gomail.ReadMessage(
			strings.NewReader(
				doc.Find("body div#mail pre").Text(),
			),
		)
		if err != nil {
			return m, err
		}
		body, err := io.ReadAll(msg.Body)
		if err != nil {
			return m, err
		}
		m = &SourceMail{
			Headers: msg.Header,
			Body:    string(body),
		}
	}
	return m, nil
}
