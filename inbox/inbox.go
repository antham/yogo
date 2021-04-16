package inbox

import (
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antham/yogo/inbox/client"
)

const itemNumber = 15

// Inbox represents a mail collection
type Inbox struct {
	Name   string `json:"name"`
	Mails  []Mail `json:"mails"`
	client client.Client
}

// NewInbox creates a new mail inbox
func NewInbox(name string) (Inbox, error) {
	client, err := client.New()
	return Inbox{
		client: client,
		Name:   name,
	}, err
}

// Get returns email at given offset
func (i *Inbox) Get(offset int) *Mail {
	if len(i.Mails) > offset {
		return &i.Mails[offset]
	}

	return nil
}

// Count returns total number of mails available in inbox
func (i *Inbox) Count() int {
	return len(i.Mails)
}

// Shrink reduces mails size to given value
func (i *Inbox) Shrink(limit int) {
	if len(i.Mails) < limit {
		return
	}

	i.Mails = i.Mails[:limit]
}

// Add appends a mail to mail list
func (i *Inbox) Add(mail Mail) {
	i.Mails = append(i.Mails, mail)
}

// Delete an email
func (i *Inbox) Delete(position int) error {
	mail := i.Mails[position]
	if err := i.client.DeleteMail(i.Name, mail.ID); err != nil {
		return err
	}

	i.Mails = append(i.Mails[:position], i.Mails[position+1:]...)
	return nil
}

// Parse retrieves all email datas
func (i *Inbox) Parse(position int) error {
	mail := &i.Mails[position]
	doc, err := i.client.GetMailPage(i.Name, mail.ID)
	if err != nil {
		return err
	}

	parseMail(doc, mail)

	return nil
}

// Flush empties an inbox
func (i *Inbox) Flush() error {
	if len(i.Mails) == 0 {
		return nil
	}

	if err := i.client.FlushMail(i.Name, i.Mails[0].ID); err != nil {
		return err
	}

	i.Mails = []Mail{}
	return nil
}

// ParseInboxPages parses inbox email in given page
func (i *Inbox) ParseInboxPages(limit int) error {
	for page := 1; page <= (limit/itemNumber)+1 && limit >= i.Count(); page++ {
		doc, err := i.client.GetMailsPage(i.Name, page)
		if err != nil {
			return err
		}

		parseInboxPage(doc, i)
		time.Sleep(1 * time.Second)
	}

	i.Shrink(limit)

	return nil
}

func parseMailID(s string) string {
	matches := regexp.MustCompile("m.php.b=.*?id=(.*)").FindStringSubmatch(s)
	if len(matches) == 2 {
		return matches[1]
	}

	return ""
}

// ParseInboxPage parses inbox email in given page
func parseInboxPage(doc *goquery.Document, inbox *Inbox) {
	doc.Find("div.um").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Find("a.lm").Attr("href")
		if !ok {
			return
		}

		var isSPAM bool
		name := s.Find("span.lmf").Text()
		mail := name

		if len(name) >= 6 && name[:6] == "[SPAM]" {
			isSPAM = true
			name = name[6:]
		}

		if strings.Contains(name, "@") {
			name = ""
		} else {
			mail = ""
		}

		if ID := parseMailID(href); ID != "" {
			mail := Mail{
				ID: ID,
				Sender: &Sender{
					Name: name,
					Mail: mail,
				},
				Title:  s.Find("span.lms").Text(),
				IsSPAM: isSPAM,
			}

			inbox.Add(mail)
		}
	})
}
