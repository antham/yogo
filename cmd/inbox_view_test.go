package cmd

import (
	"testing"
	"time"

	"github.com/antham/yogo/inbox"
	"github.com/stretchr/testify/assert"
)

func TestComputeMailOutput(t *testing.T) {
	date, err := time.Parse("2006-01-02 15:04", "2022-10-24 23:20")
	assert.NoError(t, err)

	actual, err := computeMailOutput(&inbox.Mail{ID: "test", Sender: &inbox.Sender{Mail: "test@yopmail.com"}, Title: "A title", Date: &date, Body: "test"})

	expected := `---
From  : test@yopmail.com
Title : A title
Date  : 2022-10-24 23:20
---
test
---
`
	assert.NoError(t, err)
	assert.Equal(t, expected, actual, "")
}

func TestComputeMailOutputWithAnEmptySenderAndBody(t *testing.T) {
	actual, err := computeMailOutput(&inbox.Mail{ID: "test", Sender: &inbox.Sender{}, Title: "title"})

	expected := `---
From  : <No data to display>
Title : title
Date  : <No data to display>
---
<No data to display>
---
`
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
