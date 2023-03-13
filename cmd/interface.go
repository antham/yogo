package cmd

import "github.com/antham/yogo/inbox"

type Inbox interface {
	ParseInboxPages(int) error
	Count() int
	GetMails() []inbox.Mail
	Parse(int) error
	Get(int) *inbox.Mail
	Flush() error
}
