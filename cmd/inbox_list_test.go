package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/antham/yogo/inbox"
)

type InboxMock struct {
	count                      int
	mails                      []inbox.Mail
	parseInboxPagesIntArgument int
	parseInboxPagesError       error
	parseIntArgument           int
	parseError                 error
	getIntArgument             int
	getMail                    *inbox.Mail
	flushError                 error
	deleteIntArgument          int
	deleteError                error
}

func (i *InboxMock) Count() int {
	return i.count
}

func (i *InboxMock) GetMails() []inbox.Mail {
	return i.mails
}

func (i *InboxMock) ParseInboxPages(parseInboxPagesIntArgument int) error {
	i.parseInboxPagesIntArgument = parseInboxPagesIntArgument
	return i.parseInboxPagesError
}

func (i *InboxMock) Parse(parseIntArgument int) error {
	i.parseIntArgument = parseIntArgument
	return i.parseError
}

func (i *InboxMock) Get(getIntArgument int) *inbox.Mail {
	i.getIntArgument = getIntArgument
	return i.getMail
}

func (i *InboxMock) Flush() error {
	return i.flushError
}

func (i *InboxMock) Delete(deleteIntArgument int) error {
	i.deleteIntArgument = deleteIntArgument
	return i.deleteError
}

func TestInboxList(t *testing.T) {
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
				mock.mails = []inbox.Mail{}
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
			name: "Render inbox",
			args: []string{"test", "1"},
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{}
				mock.count = 1
				mock.mails = []inbox.Mail{
					{
						ID:    "abcdefg",
						Title: "title",
						Body:  "body",
						Sender: &inbox.Sender{
							Mail: "test123@protonmail.com",
							Name: "name123",
						},
					},
				}
				return mock, nil
			},
			output: ` 1 name123 <test123@protonmail.com>
   title
`,
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
			err := inboxList(scenario.inboxBuilder)(cmd, scenario.args)
			assert.Equal(t, scenario.errExpected, err)
			assert.Equal(t, scenario.output, output.String())
			assert.Equal(t, scenario.outputErr, outputErr.String())
		})
	}
}
