package inbox

import (
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antham/yogo/inbox/internal/client"
	"github.com/antham/yogo/inbox/internal/mail"
)

const itemNumber = 15

type Mail = mail.Mail
type Sender = mail.Sender

// Inbox represents a mail collection
type Inbox struct {
	Name       string      `json:"name"`
	InboxItems []InboxItem `json:"mails"`
	client     client.Client
}

// Inbox represents a mail sumup in an inbox
type InboxItem struct {
	ID     string       `json:"id"`
	Sender *mail.Sender `json:"sender,omitempty"`
	Title  string       `json:"title"`
	Date   *time.Time   `json:"date,omitempty"`
	Body   string       `json:"body,omitempty"`
	IsSPAM bool         `json:"isSPAM"`
}

// NewInbox creates a new mail inbox
func NewInbox(name string) (*Inbox, error) {
	client, err := client.New()
	return &Inbox{
		client: client,
		Name:   name,
	}, err
}

// Fetch retrieves the full email content from the given
// inbox email offset
func (i *Inbox) Fetch(offset int) (*mail.Mail, error) {
	ID := &i.InboxItems[offset].ID
	doc, err := i.client.GetMailPage(i.Name, *ID)
	if err != nil {
		return nil, err
	}

	m := &mail.Mail{ID: *ID}
	mail.Parse(doc, m)
	return m, nil
}

// Count returns total number of mails available in inbox
func (i *Inbox) Count() int {
	return len(i.InboxItems)
}

// Shrink reduces mails size to given value
func (i *Inbox) Shrink(limit int) {
	if len(i.InboxItems) < limit {
		return
	}

	i.InboxItems = i.InboxItems[:limit]
}

// Add appends a mail to mail list
func (i *Inbox) Add(inboxItem InboxItem) {
	i.InboxItems = append(i.InboxItems, inboxItem)
}

// Delete an email
func (i *Inbox) Delete(position int) error {
	mail := i.InboxItems[position]
	if err := i.client.DeleteMail(i.Name, mail.ID); err != nil {
		return err
	}

	i.InboxItems = append(i.InboxItems[:position], i.InboxItems[position+1:]...)
	return nil
}

// Flush empties an inbox
func (i *Inbox) Flush() error {
	if len(i.InboxItems) == 0 {
		return nil
	}

	if err := i.client.FlushMail(i.Name, i.InboxItems[0].ID); err != nil {
		return err
	}

	i.InboxItems = []InboxItem{}
	return nil
}

func (i *Inbox) GetMails() []InboxItem {
	return i.InboxItems
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

// ParseInboxPage parses inbox email in given page
func parseInboxPage(doc *goquery.Document, inbox *Inbox) {
	doc.Find("div.m").Each(func(i int, s *goquery.Selection) {
		var isSPAM bool
		name := s.Find("span.lmf").Text()
		userEmail := name

		if len(name) >= 6 && name[:6] == "[SPAM]" {
			isSPAM = true
			name = name[6:]
		}

		if strings.Contains(name, "@") {
			name = ""
		} else {
			userEmail = ""
		}

		if ID, ok := s.Attr("id"); ok {
			inboxItem := InboxItem{
				ID: ID,
				Sender: &mail.Sender{
					Name: name,
					Mail: userEmail,
				},
				Title:  s.Find("div.lms").Text(),
				IsSPAM: isSPAM,
			}

			inbox.Add(inboxItem)
		}
	})
}
