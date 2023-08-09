package cmd

import "github.com/antham/yogo/inbox"

type Inbox interface {
	ParseInboxPages(int) error
	Count() int
	GetMails() []inbox.InboxItem
	Fetch(int) (*inbox.Mail, error)
	Flush() error
	Delete(int) error
}
