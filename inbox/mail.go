package inbox

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
)

// Sender defines a mail sender
type Sender struct {
	Mail string `json:"mail,omitempty"`
	Name string `json:"name,omitempty"`
}

// Mail is a message
type Mail struct {
	ID     string     `json:"id"`
	Sender *Sender    `json:"sender,omitempty"`
	Title  string     `json:"title"`
	Date   *time.Time `json:"date,omitempty"`
	Body   string     `json:"body,omitempty"`
	IsSPAM bool       `json:"isSPAM"`
}

func parseFrom(s string) (string, string) {
	re := regexp.MustCompile(`(?s)(.+?) <(.+?)>`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]), matches[2]
	}

	re = regexp.MustCompile(`<(.+?)>`)
	matches = re.FindStringSubmatch(s)
	if len(matches) == 2 {
		return "", matches[1]
	}

	return "", ""
}

func parseDate(s string) *time.Time {
	date, err := time.Parse("Monday, January 02, 2006 3:04:05 PM", s)
	if err != nil {
		return nil
	}

	return &date
}

func parseMail(doc *goquery.Document, mail *Mail) {
	mail.Sender = &Sender{}

	doc.Find("body div.fl .ellipsis").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 0:
			mail.Title = strings.TrimSpace(s.Text())
		case 1:
			mail.Sender.Name, mail.Sender.Mail = parseFrom(s.Text())
		case 2:
			mail.Date = parseDate(strings.Join(strings.Fields(s.Text()), " "))
		}
	})
	mail.Body = parseHTML(doc.Find("div#mail").Html())
}

func parseHTML(content string, err error) string {
	if err != nil {
		return ""
	}
	text, err := html2text.FromString(content)
	if err != nil {
		return ""
	}

	return text
}
