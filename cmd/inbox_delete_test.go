package cmd

import (
	"errors"
	"testing"

	"github.com/antham/yogo/inbox"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestInboxDelete(t *testing.T) {
	type scenario struct {
		name         string
		args         []string
		errExpected  error
		inboxBuilder inboxBuilder
	}

	scenarios := []scenario{
		{
			name: "No mails found",
			args: []string{"test", "1"},
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
			name:        "An error is thrown when deleting a message",
			args:        []string{"test", "1"},
			errExpected: errors.New("delete message error"),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{deleteError: errors.New("delete message error")}
				return mock, nil
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			cmd := &cobra.Command{}
			err := inboxDelete(scenario.inboxBuilder)(cmd, scenario.args)
			assert.Equal(t, scenario.errExpected, err)
		})
	}
}
