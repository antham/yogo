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

func TestRenderInboxWithAnEmptySenderEmail(t *testing.T) {
	actual := []string{}

	output = func(data string) {
		actual = append(actual, data)
	}

	successExit = func() {
		t.SkipNow()
	}

	in := inbox.Inbox{}
	in.Add(inbox.Mail{ID: "test", Sender: &inbox.Sender{Name: "name", Mail: ""}, Title: "title"})
	renderInboxMail(&in)

	assert.Regexp(t, ".*1.*name.*.*", actual[0], "Must display sender name")
	assert.Regexp(t, ".*title.*", actual[1], "Must display email title")
}

func TestRenderInboxWithAnEmptySenderName(t *testing.T) {
	actual := []string{}

	output = func(data string) {
		actual = append(actual, data)
	}

	successExit = func() {
		t.SkipNow()
	}

	in := inbox.Inbox{}
	in.Add(inbox.Mail{ID: "test", Sender: &inbox.Sender{Name: "", Mail: "test@test.com"}, Title: "title"})
	renderInboxMail(&in)

	assert.Regexp(t, ".*1.*test@test.com.*.*", actual[0], "Must display sender email")
	assert.Regexp(t, ".*title.*", actual[1], "Must display email title")
}
