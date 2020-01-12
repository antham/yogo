package inbox

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const refURL = "http://www.yopmail.com"
const jsFileURL = refURL + "/style/2.9/webmail.js"

var inboxURLs = map[string]string{
	"index": refURL + "/inbox.php?login=%v&p=%v&d=&ctrl=&scrl=&spam=true&v=2.9&r_c=&id=",
	"flush": refURL + "/inbox.php?login=%v&p=1&d=all&ctrl=%v&v=2.9&r_c=&id=",
}

var itemNumber = 15

// Inbox represents a mail collection
type Inbox struct {
	identifier string
	mails      []Mail
}

// Get return email at given offset
func (i *Inbox) Get(offset int) *Mail {
	if len(i.mails) > offset {
		return &i.mails[offset]
	}

	return nil
}

// Count return total number of mails available in inbox
func (i *Inbox) Count() int {
	return len(i.mails)
}

// Shrink reduce mails size to given value
func (i *Inbox) Shrink(limit int) {
	if len(i.mails) < limit {
		return
	}

	i.mails = i.mails[:limit]
}

// GetAll return all emails
func (i *Inbox) GetAll() []Mail {
	return i.mails
}

// GetIdentifier return mailbox name
func (i *Inbox) GetIdentifier() string {
	return i.identifier
}

// Add append a mail to mails
func (i *Inbox) Add(mail Mail) {
	i.mails = append(i.mails, mail)
}

// Delete an email
func (i *Inbox) Delete(position int) error {
	mail := i.mails[position]
	if err := send(fmt.Sprintf(mailURLs["delete"], i.GetIdentifier(), strings.TrimLeft(mail.ID, "m"))); err != nil {
		return err
	}

	i.mails = append(i.mails[:position], i.mails[position+1:]...)
	return nil
}

// Parse retrieve all email datas
func (i *Inbox) Parse(position int) error {
	mail := &i.mails[position]
	URL := fmt.Sprintf(mailURLs["get"], i.identifier, mail.ID)

	r, err := buildReader("GET", URL, map[string]string{"Cookie": fmt.Sprintf("compte=%s", i.identifier)}, nil)
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
	if len(i.mails) == 0 {
		return nil
	}

	if err := send(fmt.Sprintf(inboxURLs["flush"], i.identifier, strings.TrimLeft(i.mails[0].ID, "m"))); err != nil {
		return err
	}

	i.mails = []Mail{}
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
	inbox := Inbox{identifier: identifier}

	for page := 1; page <= (limit/itemNumber)+1 && limit >= inbox.Count(); page++ {
		URL, err := urlDecorator(fmt.Sprintf(inboxURLs["index"], identifier, page))
		if err != nil {
			return nil, err
		}

		r, err := buildReader("GET", URL, map[string]string{"Cookie": fmt.Sprintf("compte=%s", identifier)}, nil)
		if err != nil {
			return nil, err
		}

		doc, err := fetchFromReader(r)
		if err != nil {
			return nil, err
		}

		parseInboxPage(doc, &inbox)
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

		if ID := parseMailID(href); ID != "" {
			mail := Mail{
				ID:    ID,
				Title: s.Find("span.lmf").Text(),
				SumUp: s.Find("span.lms").Text(),
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

	doc, err = fetchURL(jsFileURL)
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
