package cmd

import (
	"fmt"
	"strconv"
)

func parseMailAndOffsetArgs(args []string) (string, int) {
	if len(args) != 2 {
		perror(fmt.Errorf("An inbox name without @yopmail.com and an offset are required"))

		errorExit()
	}

	offset, err := strconv.Atoi(args[1])

	if err != nil {
		perror(fmt.Errorf(`argument "%s" must be an integer`, args[1]))

		errorExit()
	}

	if offset < 1 {
		perror(fmt.Errorf(`argument "%d" must be greater than 0`, offset))

		errorExit()
	}

	return args[0], offset
}

func checkOffset(count int, offset int) {
	if count == 0 {
		perror(fmt.Errorf("Inbox is empty"))

		errorExit()
	}

	if count < offset-1 {
		perror(fmt.Errorf("Lower your offset value"))

		errorExit()
	}
}
