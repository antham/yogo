package inbox

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var mailURLs = map[string]string{
	"get":    refURL + "/m.php?b=%v&id=%v",
	"delete": refURL + "/inbox.php?login=%v&p=1&d=%v&ctrl=&scrl=0&spam=true&v=2.9&r_c=",
}

// Sender defines a mail sender
type Sender struct {
	Mail string
	Name string
}

// Mail is a message
type Mail struct {
	ID      string
	Sender  Sender
	SumUp   string
	Title   string
	Date    time.Time
	Headers []string
	Body    string
}

func parseFrom(s string) (string, string) {
	re := regexp.MustCompile(`.*?:\s*"?(.*?)"?\s*<(.*?)>`)
	matches := re.FindStringSubmatch(s)

	if len(matches) == 3 {
		return matches[1], matches[2]
	}

	return "", ""
}

func parseDate(s string) time.Time {
	re := regexp.MustCompile(`.*?(\d+/\d+/\d+).*?(\d+:\d+)`)
	matches := re.FindStringSubmatch(s)

	if len(matches) != 3 {
		return time.Time{}
	}

	date, err := time.Parse("02/01/2006 15:04", fmt.Sprintf("%v %v", matches[1], matches[2]))
	if err != nil {
		return time.Time{}
	}

	return date
}

func parseMail(doc *goquery.Document, mail *Mail) {
	doc.Find("body").Each(func(i int, s *goquery.Selection) {
		mail.Sender = Sender{}
		mail.Sender.Name, mail.Sender.Mail = parseFrom(s.Find("div#mailhaut div:nth-child(2)").Text())
		mail.Date = parseDate(s.Find("div#mailhaut div:nth-child(4)").Text())

		mail.Body = parseHTML(s.Find("div#mailmillieu").Html())
		mail.Title = strings.TrimSpace(s.Find("div#mailhaut .f16").Text())
	})
}
