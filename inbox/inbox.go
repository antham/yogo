package inbox

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func getJsFileURL() string {
	return refURL + "/style/" + apiVersion + "/webmail.js"
}

func getInboxURLs(key string) string {
	return map[string]string{
		"index": refURL + "/inbox.php?login=%v&p=%v&d=&ctrl=&scrl=&spam=true&v=" + apiVersion + "&r_c=&id=",
		"flush": refURL + "/inbox.php?login=%v&p=1&d=all&ctrl=%v&v=" + apiVersion + "&r_c=&id=none&scrl=&spam=true",
	}[key]
}

const itemNumber = 15

// Inbox represents a mail collection
type Inbox struct {
	Name  string `json:"name"`
	Mails []Mail `json:"mails"`
}

// Get return email at given offset
func (i *Inbox) Get(offset int) *Mail {
	if len(i.Mails) > offset {
		return &i.Mails[offset]
	}

	return nil
}

// Count return total number of mails available in inbox
func (i *Inbox) Count() int {
	return len(i.Mails)
}

// Shrink reduce mails size to given value
func (i *Inbox) Shrink(limit int) {
	if len(i.Mails) < limit {
		return
	}

	i.Mails = i.Mails[:limit]
}

// Add append a mail to mails
func (i *Inbox) Add(mail Mail) {
	i.Mails = append(i.Mails, mail)
}

// Delete an email
func (i *Inbox) Delete(position int) error {
	mail := i.Mails[position]
	URL, err := urlDecorator(fmt.Sprintf(getMailURLs("delete"), i.Name, strings.TrimLeft(mail.ID, "m")))
	if err != nil {
		return err
	}
	if err := send(URL, createCompteCookie(i.Name)); err != nil {
		return err
	}

	i.Mails = append(i.Mails[:position], i.Mails[position+1:]...)
	return nil
}

// Parse retrieve all email datas
func (i *Inbox) Parse(position int) error {
	mail := &i.Mails[position]
	URL := fmt.Sprintf(getMailURLs("get"), i.Name, mail.ID)

	r, err := buildReader("GET", URL, createCompteCookie(i.Name), nil)
	if err != nil {
		return err
	}

	doc, err := fetchFromReader(r)
	if err != nil {
		return err
	}

	parseMail(doc, mail)

	return nil
}

// Flush empty an inbox
func (i *Inbox) Flush() error {
	if len(i.Mails) == 0 {
		return nil
	}

	URL, err := urlDecorator(fmt.Sprintf(getInboxURLs("flush"), i.Name, strings.TrimLeft(i.Mails[0].ID, "m")))
	if err != nil {
		return err
	}

	if err := send(URL, createCompteCookie(i.Name)); err != nil {
		return err
	}

	i.Mails = []Mail{}
	return nil
}

func parseMailID(s string) string {
	matches := regexp.MustCompile("m.php.b=.*?id=(.*)").FindStringSubmatch(s)
	if len(matches) == 2 {
		return matches[1]
	}

	return ""
}

// ParseInboxPages parse inbox email in given page
func ParseInboxPages(identifier string, limit int) (*Inbox, error) {
	inbox := Inbox{Name: identifier}

	for page := 1; page <= (limit/itemNumber)+1 && limit >= inbox.Count(); page++ {
		URL, err := urlDecorator(fmt.Sprintf(getInboxURLs("index"), identifier, page))
		if err != nil {
			return nil, err
		}

		r, err := buildReader("GET", URL, createCompteCookie(identifier), nil)
		if err != nil {
			return nil, err
		}

		doc, err := fetchFromReader(r)
		if err != nil {
			return nil, err
		}

		parseInboxPage(doc, &inbox)
		time.Sleep(1 * time.Second)
	}

	inbox.Shrink(limit)

	return &inbox, nil
}

// ParseInboxPage parse inbox email in given page
func parseInboxPage(doc *goquery.Document, inbox *Inbox) {
	doc.Find("div.um").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find("a.lm").Attr("href")
		if !ok {
			return
		}

		var isSPAM bool
		title := s.Find("span.lmf").Text()

		if len(title) >= 6 && title[:6] == "[SPAM]" {
			isSPAM = true
			title = title[6:]
		}

		if ID := parseMailID(href); ID != "" {
			mail := Mail{
				ID:    ID,
				Title: title,
				SumUp: func() *string {
					v := s.Find("span.lms").Text()
					return &v
				}(),
				IsSPAM: isSPAM,
			}

			inbox.Add(mail)
		}
	})
}

func urlDecorator(URL string) (string, error) {
	doc, err := fetchURL(refURL)
	if err != nil {
		return "", err
	}

	var yp string

	doc.Find("#yp").Each(func(i int, s *goquery.Selection) {
		var ok bool
		yp, ok = s.Attr("value")
		if !ok {
			err = errors.New("no attribute yp found")
		}
	})

	if err != nil {
		return "", err
	}

	doc, err = fetchURL(getJsFileURL())
	if err != nil {
		return "", err
	}

	yj := regexp.MustCompile("&yj=(.*?)&").FindStringSubmatch(doc.Text())[1]

	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	q := u.Query()
	q.Add("yp", yp)
	q.Add("yj", yj)

	return (&url.URL{
		Scheme:   u.Scheme,
		Host:     u.Host,
		Path:     u.Path,
		RawQuery: q.Encode(),
	}).String(), nil
}
