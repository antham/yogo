package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/antham/yogo/v4/internal/inbox"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestInboxFlush(t *testing.T) {
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
			name: "No mails found",
			args: []string{"test", "1"},
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{}
				mock.items = []inbox.InboxItem{}
				return mock, nil
			},
			output: `Inbox "test" successfully flushed
`,
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
			name:        "An error is thrown when flushing inbox",
			args:        []string{"test", "1"},
			errExpected: errors.New("flush inbox error"),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{flushError: errors.New("flush inbox error")}
				mock.count = 1
				mock.items = []inbox.InboxItem{
					{
						ID:      "abcdefg",
						Subject: "subject",
						Body:    "body",
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
			name: "Inbox flushed successfully",
			args: []string{"test", "1"},
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{}
				mock.count = 1
				mock.items = []inbox.InboxItem{
					{
						ID:      "abcdefg",
						Subject: "subject",
						Body:    "body",
						Sender: &inbox.Sender{
							Mail: "test123",
							Name: "name123",
						},
					},
				}
				return mock, nil
			},
			output: `Inbox "test" successfully flushed
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
			err := inboxFlush(scenario.inboxBuilder)(cmd, scenario.args)
			assert.Equal(t, scenario.errExpected, err)
			assert.Equal(t, scenario.output, output.String())
			assert.Equal(t, scenario.outputErr, outputErr.String())
		})
	}
}
