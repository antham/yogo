package cmd

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/antham/yogo/inbox"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestInboxShow(t *testing.T) {
	type scenario struct {
		name         string
		args         []string
		errExpected  error
		inboxBuilder inboxBuilder
		output       string
		outputErr    string
	}

	scenarios := []scenario{
		{
			name:        "No mails found",
			args:        []string{"test", "1"},
			errExpected: errors.New("inbox is empty"),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{}
				mock.items = []inbox.InboxItem{}
				return mock, nil
			},
		},
		{
			name:        "Failure when parsing offset",
			args:        []string{"test", "-1"},
			errExpected: errors.New(`offset "-1" must be greater than 0`),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{}
				mock.items = []inbox.InboxItem{}
				return mock, nil
			},
		},
		{
			name:        "An error is thrown in inbox builder",
			args:        []string{"test", "1"},
			errExpected: errors.New("inbox builder error"),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{}
				return mock, errors.New("inbox builder error")
			},
		},
		{
			name:        "An error is thrown in parse inbox pages",
			args:        []string{"test", "1"},
			errExpected: errors.New("inbox pages error"),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{parseInboxPagesError: errors.New("inbox pages error")}
				return mock, nil
			},
		},
		{
			name:        "An error is thrown when parsing mail",
			args:        []string{"test", "1"},
			errExpected: errors.New("parse email error"),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{parseInboxPagesError: errors.New("parse email error")}
				mock.count = 1
				mock.items = []inbox.InboxItem{
					{
						ID:    "abcdefg",
						Title: "title",
						Body:  "body",
						Sender: &inbox.Sender{
							Mail: "test123",
							Name: "name123",
						},
					},
				}
				return mock, nil
			},
		},
		{
			name: "Offset to high compared to the number of emails",
			args: []string{"test", "2"},
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{fetchMail: nil}
				mock.count = 1
				mock.items = []inbox.InboxItem{
					{
						ID:    "abcdefg",
						Title: "title",
						Body:  "body",
						Sender: &inbox.Sender{
							Mail: "test123",
							Name: "name123",
						},
					},
				}
				return mock, nil
			},
		},
		{
			name:        "No mail found",
			args:        []string{"test", "1"},
			errExpected: errors.New("inbox is empty"),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{fetchMail: nil}
				return mock, nil
			},
		},
		{
			name: "Output the mail",
			args: []string{"test", "1"},
			output: `---
From  : name123 <test123>
Title : title
Date  : 2001-01-01 00:00
---
body
---

`,
			inboxBuilder: func(name string) (Inbox, error) {
				now, _ := time.Parse("2006-01-02", "2001-01-01")
				mock := &InboxMock{}
				mock.count = 1
				mock.items = []inbox.InboxItem{
					{
						ID:    "abcdefg",
						Title: "title",
						Body:  "body",
						Sender: &inbox.Sender{
							Mail: "test123",
							Name: "name123",
						},
						Date: &now,
					},
				}
				mock.fetchMail = &inbox.Mail{
					ID:    "abcdefg",
					Title: "title",
					Body:  "body",
					Sender: &inbox.Sender{
						Mail: "test123",
						Name: "name123",
					},
					Date: &now,
				}
				return mock, nil
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			var output bytes.Buffer
			var outputErr bytes.Buffer
			cmd := &cobra.Command{}
			cmd.SetOut(&output)
			cmd.SetErr(&outputErr)
			err := inboxShow(scenario.inboxBuilder)(cmd, scenario.args)
			assert.Equal(t, scenario.errExpected, err)
			assert.Equal(t, scenario.output, output.String())
			assert.Equal(t, scenario.outputErr, outputErr.String())
		})
	}
}
