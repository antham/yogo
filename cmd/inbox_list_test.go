package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/yogo/inbox"
)

func TestRenderInboxMailWithEmptyInbox(t *testing.T) {
	info = func(msg string) {
		assert.Equal(t, "Inbox is empty", msg, "Must return an info message")
	}

	successExit = func() {
		t.SkipNow()
	}

	in := inbox.Inbox{}

	renderInboxMail(&in)
}

func TestRenderInbox(t *testing.T) {
	actual := []string{}

	output = func(data string) {
		actual = append(actual, data)
	}

	successExit = func() {
		t.SkipNow()
	}

	in := inbox.Inbox{}
	in.Add(inbox.Mail{ID: "test", Title: "title", SumUp: "Sum up"})
	renderInboxMail(&in)

	assert.Equal(t, " \x1b[32m1\x1b[0m \x1b[33mtitle\x1b[0m\n", actual[0], "Must display email title")
	assert.Equal(t, " \x1b[36mSum up\x1b[0m\n\n", actual[1], "Must display email sum up")
}
