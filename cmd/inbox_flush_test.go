package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInboxFlushWithNoArguments(t *testing.T) {

	perror = func(err error) {
		assert.EqualError(t, err, "One argument mandatory", "Must return an error")
	}

	errorExit = func() {
		t.SkipNow()
	}

	os.Args = []string{"", "inbox", "flush"}

	assert.NoError(t, RootCmd.Execute())
}
