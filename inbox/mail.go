package inbox

import (
	"fmt"
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
	re := regexp.MustCompile(`.+?:\s*"?(.+?)"?\s*<(.+?)>`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 3 {
		return strings.TrimSpace(matches[1]), matches[2]
	}

	re = regexp.MustCompile(`.+?:\s*(.+)`)
	matches = re.FindStringSubmatch(s)
	if len(matches) == 2 {
		return "", matches[1]
	}

	return "", ""
}

func parseDate(s string) *time.Time {
	re := regexp.MustCompile(`.*?(\d+/\d+/\d+).*?(\d+:\d+)`)
	matches := re.FindStringSubmatch(s)

	if len(matches) != 3 {
		return nil
	}

	date, err := time.Parse("02/01/2006 15:04", fmt.Sprintf("%v %v", matches[1], matches[2]))
	if err != nil {
		return nil
	}

	return &date
}

func parseMail(doc *goquery.Document, mail *Mail) {
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		mail.Sender = &Sender{}
		mail.Sender.Name, mail.Sender.Mail = parseFrom(s.Find("div#mailhaut div:nth-child(2)").Text())
		mail.Date = parseDate(s.Find("div#mailhaut div:nth-child(4)").Text())

		mail.Body = parseHTML(s.Find("div#mailmillieu").Html())
		mail.Title = strings.TrimSpace(s.Find("div#mailhaut .f16").Text())
	})
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
