package cmd

import (
	"github.com/antham/yogo/v4/internal/inbox"
)

type Inbox interface {
	inbox.Render
	ParseInboxPages(int) error
	Count() int
	GetMails() []inbox.InboxItem
	Fetch(int) (inbox.Render, error)
	Flush() error
	Delete(int) error
}
