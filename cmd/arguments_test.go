package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMailAndOffsetArgsWithNoArguments(t *testing.T) {
	perror = func(err error) {
		assert.EqualError(t, err, "An inbox name without @yopmail.com and an offset are required", "Must return an error")
	}

	errorExit = func() {
		t.SkipNow()
	}

	parseMailAndOffsetArgs([]string{})
}

func TestParseMailAndOffsetArgsWithSecondArgumentAString(t *testing.T) {
	perror = func(err error) {
		assert.EqualError(t, err, `argument "test" must be an integer`, "Must return an error")
	}

	errorExit = func() {
		t.SkipNow()
	}

	parseMailAndOffsetArgs([]string{"test", "test"})
}

func TestParseMailAndOffsetArgsWithSecondArgumentLessThan0(t *testing.T) {
	perror = func(err error) {
		assert.EqualError(t, err, `argument "0" must be greater than 0`, "Must return an error")
	}

	errorExit = func() {
		t.SkipNow()
	}

	parseMailAndOffsetArgs([]string{"test", "0"})
}
