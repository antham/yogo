package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOffset(t *testing.T) {
	type scenario struct {
		name           string
		err            error
		offsetArg      string
		offsetExpected int
	}

	scenarios := []scenario{
		{
			name:      "offset is a string",
			offsetArg: "test",
			err:       errors.New(`offset "test" must be an integer`),
		},
		{
			name:      "offset lower than 0",
			offsetArg: "0",
			err:       errors.New(`offset "0" must be greater than 0`),
		},
		{
			name:           "valid offset",
			offsetArg:      "1",
			offsetExpected: 1,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			offset, err := parseOffset(scenario.offsetArg)
			if scenario.err != nil {
				assert.EqualError(t, err, scenario.err.Error())
			} else {
				assert.Equal(t, scenario.offsetExpected, offset)
			}
		})
	}
}

func TestNormalizeInboxName(t *testing.T) {
	type scenario struct {
		name          string
		inboxArg      string
		inboxExpected string
	}

	scenarios := []scenario{
		{
			name:          "uppercased inbox",
			inboxArg:      "TeSt",
			inboxExpected: "test",
		},
		{
			name:          "full email provided",
			inboxArg:      "test@yopmail.com",
			inboxExpected: "test",
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			inbox := normalizeInboxName(scenario.inboxArg)
			assert.Equal(t, scenario.inboxExpected, inbox)
		})
	}
}

func TestCheckOffset(t *testing.T) {
	type scenario struct {
		name string
		args []int
		err  error
	}

	scenarios := []scenario{
		{
			name: "second argument as string",
			args: []int{1, 3},
			err:  errors.New(`lower your offset value`),
		},
		{
			name: "regular offset",
			args: []int{0, 1},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			err := checkOffset(scenario.args[0], scenario.args[1])
			if scenario.err != nil {
				assert.EqualError(t, err, scenario.err.Error())
			}
		})
	}
}
