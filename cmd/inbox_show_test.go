package cmd

import (
	"errors"
	"testing"

	"github.com/antham/yogo/inbox"
	"github.com/stretchr/testify/assert"
)

func TestInboxShow(t *testing.T) {
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
			name:        "An error is thrown when parsing mail",
			args:        []string{"test", "1"},
			errExpected: errors.New("parse email error"),
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{parseError: errors.New("parse email error")}
				return mock, nil
			},
		},
		{
			name: "No mail found",
			args: []string{"test", "1"},
			inboxBuilder: func(name string) (Inbox, error) {
				mock := &InboxMock{getMail: nil}
				return mock, nil
			},
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()
			err := inboxShow(scenario.inboxBuilder)(nil, scenario.args)
			assert.Equal(t, scenario.errExpected, err)
		})
	}
}
