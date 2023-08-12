package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInbox(t *testing.T) {
	in, err := newInbox("test")
	assert.NoError(t, err)
	assert.NotEmpty(t, in)
}
