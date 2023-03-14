package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseMailAndOffsetArgs(args []string) (string, int, error) {
	offset, err := strconv.Atoi(args[1])
	if err != nil {
		return "", 0, fmt.Errorf(`argument "%s" must be an integer`, args[1])
	}

	if offset < 1 {
		return "", 0, fmt.Errorf(`argument "%d" must be greater than 0`, offset)
	}

	// Providing an uppercased email triggers a panic.
	// In the web interface there is a redirection to
	// the inbox with the address lowercased so we mimic
	// this behaviour
	return strings.ToLower(args[0]), offset, nil
}

func checkOffset(count int, offset int) error {
	if count < offset-1 {
		return errors.New("Lower your offset value")
	}

	if count == 0 {
		info("Inbox is empty")
	}
	return nil
}
