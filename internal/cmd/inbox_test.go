package cmd

import (
	"testing"

	"github.com/antham/yogo/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNewInbox(t *testing.T) {
	in, err := newInbox[client.MailHTMLDoc]("test")
	assert.NoError(t, err)
	assert.NotEmpty(t, in)
}
