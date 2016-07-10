package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/yogo/inbox"
)

func TestInboxListWithNoArguments(t *testing.T) {

	perror = func(err error) {
		assert.EqualError(t, err, "Two arguments mandatory", "Must return an error")
	}

	errorExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "inbox", "list"}

	RootCmd.Execute()
}

func TestInboxListWithOneArgument(t *testing.T) {

	perror = func(err error) {
		assert.EqualError(t, err, "Two arguments mandatory", "Must return an error")
	}

	errorExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "inbox", "list", "test"}

	RootCmd.Execute()
}

func TestInboxListWithStringAsLastArgument(t *testing.T) {

	perror = func(err error) {
		assert.EqualError(t, err, `"test" must be an integer`, "Must return an error")
	}

	errorExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "inbox", "list", "test", "test"}

	RootCmd.Execute()
}

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

	assert.Equal(t, " 1 title\n", actual[0], "Must display email title")
	assert.Equal(t, " Sum up\n\n", actual[1], "Must display email sum up")
}
