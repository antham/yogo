package mail

import "time"
import "github.com/PuerkitoBio/goquery"
import "fmt"
import "log"
import "regexp"
import "net/http"

import "strings"

var getUrl = "http://www.yopmail.com/mail.php?b=%v&id=%v"
var deleteUrl = "http://www.yopmail.com/inbox.php?login=%v&p=1&d=%v&ctrl=&scrl=0&spam=true&v=2.6&r_c="

type Mail struct {
	id         string
	mail       string
	FromString string
	FromMail   string
	SumUp      string
	Title      string
	Date       time.Time
	headers    []string
	Body       string
}

func NewMail(mail string, id string, title string, sumUp string) *Mail {
	return &Mail{
		id:    id,
		mail:  mail,
		Title: title,
		SumUp: sumUp,
	}
}

func (m *Mail) Fetch() {

	doc, err := goquery.NewDocument(fmt.Sprintf(getUrl, m.mail, m.id))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("body").Each(func(i int, s *goquery.Selection) {

		fromString, fromMail := func(s *goquery.Selection) (string, string) {

			re := regexp.MustCompile(".*?: (.*?)<(.*?)>")

			matches := re.FindStringSubmatch(s.Find("div#mailhaut div:nth-child(2)").Text())

			if len(matches) == 3 {
				return matches[1], matches[2]
			}

			return "", ""
		}(s)

		date := func(s *goquery.Selection) time.Time {
			re := regexp.MustCompile(".*?(\\d+/\\d+/\\d+).*?(\\d+:\\d+)")

			matches := re.FindStringSubmatch(s.Find("div#mailhaut div:nth-child(4)").Text())

			if len(matches) != 3 {
				return time.Time{}
			}

			date, error := time.Parse("02/01/2006 15:04", fmt.Sprintf("%v %v", matches[1], matches[2]))

			if error != nil {
				return time.Time{}
			}

			return date
		}(s)

		m.FromString = fromString
		m.FromMail = fromMail
		m.Date = date
		m.Body = strings.TrimSpace(s.Find("div#mailmillieu").Text())
		m.Title = strings.TrimSpace(s.Find("div#mailhaut .f16").Text())
	})
}

func (m *Mail) Delete() {
	http.Get(fmt.Sprintf(deleteUrl, m.mail, strings.TrimLeft(m.id, "m")))
}
