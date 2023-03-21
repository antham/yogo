package cmd

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMailAndOffset(t *testing.T) {
	type scenario struct {
		name   string
		args   []string
		err    error
		offset int
		inbox  string
	}

	scenarios := []scenario{
		{
			name: "second argument as string",
			args: []string{"test", "test"},
			err:  errors.New(`offset "test" must be an integer`),
		},
		{
			name: "offset lower than 0",
			args: []string{"test", "0"},
			err:  errors.New(`offset "0" must be greater than 0`),
		},
		{
			name:   "regular inbox",
			args:   []string{"test", "1"},
			offset: 1,
			inbox:  "test",
		},
		{
			name:   "uppercased inbox",
			args:   []string{"TeSt", "1"},
			offset: 1,
			inbox:  "test",
		},
		{
			name:   "full email provided",
			args:   []string{"test@yopmail.com", "1"},
			offset: 1,
			inbox:  "test",
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			t.Parallel()

			inbox, offset, err := parseMailAndOffsetArgs(scenario.args)
			if scenario.err != nil {
				assert.EqualError(t, err, scenario.err.Error())
			} else {
				assert.Equal(t, scenario.offset, offset)
				assert.Equal(t, scenario.inbox, inbox)
			}
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
