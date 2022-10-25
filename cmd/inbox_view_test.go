package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/antham/yogo/inbox"
	"github.com/stretchr/testify/assert"
)

func TestRenderMail(t *testing.T) {
	actual := []string{}

	output = func(data string) {
		actual = append(actual, data)
	}

	successExit = func() {
		t.SkipNow()
	}

	date, err := time.Parse("2006-01-02 15:04", "2022-10-24 23:20")
	assert.NoError(t, err)

	renderMail(&inbox.Mail{ID: "test", Sender: &inbox.Sender{Mail: "test@yopmail.com"}, Title: "A title", Date: &date, Body: "test"})

	expected := `---
From  : test@yopmail.com
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`
	assert.Equal(t, expected, strings.Join(actual, ""))
}

func TestRenderMailWithAnEmptySenderAndBody(t *testing.T) {
	actual := []string{}

	output = func(data string) {
		actual = append(actual, data)
	}

	successExit = func() {
		t.SkipNow()
	}

	renderMail(&inbox.Mail{ID: "test", Sender: &inbox.Sender{}, Title: "title"})

	expected := `---
From  : No data to display
Title : title
Date  : No data to display
---
No data to display
---
`
	assert.Equal(t, expected, strings.Join(actual, ""))
}
