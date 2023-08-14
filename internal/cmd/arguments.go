package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func parseOffset(offset string) (int, error) {
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return 0, fmt.Errorf(`offset "%s" must be an integer`, offset)
	}

	if offsetInt < 1 {
		return 0, fmt.Errorf(`offset "%d" must be greater than 0`, offsetInt)
	}

	// Providing an uppercased email triggers a panic.
	// In the web interface there is a redirection to
	// the inbox with the address lowercased so we mimic
	// this behaviour
	return offsetInt, nil
}

func normalizeInboxName(inboxName string) string {
	// Providing an uppercased email triggers a panic.
	// In the web interface there is a redirection to
	// the inbox with the address lowercased so we mimic
	// this behaviour
	return strings.ReplaceAll(strings.ToLower(inboxName), "@yopmail.com", "")
}

func checkOffset(count int, offset int) error {
	if count < offset-1 {
		return errors.New("lower your offset value")
	}

	if count == 0 {
		return errors.New("inbox is empty")
	}
	return nil
}
